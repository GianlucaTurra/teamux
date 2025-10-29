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

// TODO: remove these methods and call tmux directly

func (s Session) Open() error {
	return tmux.CreateSession(s.Name, s.WorkingDirectory)
}

func (s Session) IsOpen() bool {
	return tmux.IsSessionOpen(s.Name)
}

func (s Session) Close() error {
	return tmux.KillSession(s.Name)
}

func (s Session) Switch() error {
	return tmux.SwitchToSession(s.Name)
}
