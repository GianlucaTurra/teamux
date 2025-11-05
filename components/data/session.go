// Package data declares the data structure to map db entities and the
// functions to interact with them
package data

import (
	"github.com/GianlucaTurra/teamux/tmux"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Name             string
	WorkingDirectory string
	Windows          []Window `gorm:"many2many:session_windows;"`
}

func ReadAllSessions(db *gorm.DB) ([]Session, error) {
	var sessions []Session
	err := db.Model(&Session{}).Preload("Windows").Find(&sessions).Error
	return sessions, err
}

func CreateSession(name string, workingDirectory string, connector Connector) (int, error) {
	session := Session{Name: name, WorkingDirectory: workingDirectory}
	result := gorm.WithResult()
	err := gorm.G[Session](connector.DB, result).Create(connector.Ctx, &session)
	return int(result.RowsAffected), err
}

func (s Session) Save(connector Connector) (int, error) {
	return gorm.G[Session](connector.DB).Updates(connector.Ctx, s)
}

func (s Session) Delete(connector Connector) (int, error) {
	return gorm.G[Session](connector.DB).Where("id = ?", s.ID).Delete(connector.Ctx)
}

func (s Session) Open() error {
	if err := tmux.NewSession(s.Name, s.WorkingDirectory); err != nil {
		return err
	}
	// TODO: is it better to stop the process at the first error or to load
	// everything that's ok and report the error?
	for _, window := range s.Windows {
		if err := window.OpenWithTarget(s.Name); err != nil {
			return err
		}
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
