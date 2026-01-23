// Package tests
package tests

import (
	"context"
	"log"

	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup() database.Connector {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(
		&sessions.Session{},
		&windows.Window{},
		&panes.Pane{},
	)
	return database.Connector{DB: db, Ctx: context.Background()}
}
