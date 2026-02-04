package tests

import (
	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	"gorm.io/gorm"
)

func MigrateTestDB(db *gorm.DB) error {
	err := db.AutoMigrate(
		&sessions.Session{},
		&windows.Window{},
		&panes.Pane{},
	)
	return err
}
