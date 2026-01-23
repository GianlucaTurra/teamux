package sessions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/database"
	"github.com/GianlucaTurra/teamux/tmux"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Name             string
	WorkingDirectory string
	Windows          []windows.Window `gorm:"many2many:session_windows;"`
}

func ReadAllSessions(db *gorm.DB) ([]Session, error) {
	var sessions []Session
	err := db.Model(&Session{}).Preload("Windows").Preload("Windows.Panes").Find(&sessions).Error
	return sessions, err
}

func CreateSession(name string, workingDirectory string, connector database.Connector) (int, error) {
	if strings.TrimSpace(name) == "" {
		err := errors.New("session name cannot be empty")
		return 0, err
	}
	if strings.TrimSpace(workingDirectory) == "" {
		workingDirectory = "~/"
	}
	session := Session{Name: name, WorkingDirectory: workingDirectory}
	result := gorm.WithResult()
	err := gorm.G[Session](connector.DB, result).Create(connector.Ctx, &session)
	return int(result.RowsAffected), err
}

func (s Session) Save(connector database.Connector) (int, error) {
	return gorm.G[Session](connector.DB).Updates(connector.Ctx, s)
}

func (s Session) Delete(connector database.Connector) (int, error) {
	return gorm.G[Session](connector.DB).Where("id = ?", s.ID).Delete(connector.Ctx)
}

// Open creates the new session and cascades to all Windows and all their panes
// if a window has no WorkingDirectory it is set to the WorkingDirectory of the
// session.
// Opening the session this way creates an empty window at first, to avoid
// confusion it is deleted and following Windows are reordered
func (s Session) Open() error {
	if err := tmux.NewSession(s.Name, s.WorkingDirectory); err != nil {
		return err
	}
	// TODO: is it better to stop the process at the first error or to load
	// everything that's ok and report the error?
	for _, window := range s.Windows {
		if strings.TrimSpace(window.WorkingDirectory) == "" {
			window.WorkingDirectory = s.WorkingDirectory
		}
		if err := window.OpenWithTarget(s.Name); err != nil {
			return err
		}
	}
	// TODO: handle error type
	if err := tmux.KillWindow(s.Name + ":0"); err != nil {
		return tmux.NewWarning(fmt.Sprintf("unable to kill empty window of %s", s.Name))
	}
	// TODO: handle error type
	if err := tmux.ReorderWindows(s.Name); err != nil {
		return tmux.NewWarning(fmt.Sprintf("unable to reoder windows in %s", s.Name))
	}
	return nil
}

func (s Session) IsOpen() bool {
	return tmux.HasSession(s.Name)
}

func (s Session) Close() error {
	return tmux.KillSession(s.Name)
}

func (s Session) Switch() error {
	return tmux.SwitchToSession(s.Name)
}
