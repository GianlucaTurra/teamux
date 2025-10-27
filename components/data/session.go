// Package data declares the data structure to map db entities and the
// functions to interact with them
package data

import (
	"fmt"
	"os/exec"

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

// TODO: the following methods should be in a proper package

// Open translates the session object to a tmux command to open a new session.
func (s Session) Open() error {
	newSessionCmd := fmt.Sprintf("tmux new-session -d -s \"%s\" -c %s", s.Name, s.WorkingDirectory)
	cmd := exec.Command("sh", "-c", newSessionCmd)
	return cmd.Run()
}

func (s Session) IsOpen() bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux has-session -t \"%s\"", s.Name))
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (s Session) Close() error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux kill-session -t \"%s\"", s.Name))
	return cmd.Run()
}

func (s Session) Switch() error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux switch -t \"%s\"", s.Name))
	return cmd.Run()
}
