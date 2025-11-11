package common

import (
	"context"

	"github.com/GianlucaTurra/teamux/components/data"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetProdDB() (*gorm.DB, error) {
	return getDB("teamux.db")
}

func GetTestDB() *data.Connector {
	db, err := getDB("file::memory:?cache=shared")
	if err != nil {
		return nil
	}
	if err := db.AutoMigrate(
		&data.Session{},
		&data.Window{},
		&data.Pane{},
	); err != nil {
		return nil
	}
	return &data.Connector{DB: db, Ctx: context.Background()}
}

func getDB(name string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(name), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
}
