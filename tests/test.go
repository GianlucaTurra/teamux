// Package tests
package tests

import (
	"context"
	"log"

	"github.com/GianlucaTurra/teamux/components/data"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup() data.Connector {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(
		&data.Session{},
		&data.Window{},
		&data.Pane{},
	)
	return data.Connector{DB: db, Ctx: context.Background()}
}
