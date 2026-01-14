package common

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GianlucaTurra/teamux/components/data"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	ProdDB = "teamux/teamux.db"
	TestDB = "file::memory:?cache=shared"
)

func GetProdDB() (*gorm.DB, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error detecting user config directory: %v", err)
	}
	prodDBPath := fmt.Sprintf("%s/%s", userConfigDir, ProdDB)
	return getDB(prodDBPath)
}

func GetTestDB() *data.Connector {
	db, err := getDB(TestDB)
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
