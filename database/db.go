package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connector struct {
	DB  *gorm.DB
	Ctx context.Context
}

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

// FIXME: this causes cycle imports. Do it in the test package
func GetTestDB() *Connector {
	return nil
	// db, err := getDB(TestDB)
	// if err != nil {
	// 	return nil
	// }
	// if err := db.AutoMigrate(
	// 	&sessions.Session{},
	// 	&windows.Window{},
	// 	&panes.Pane{},
	// ); err != nil {
	// 	return nil
	// }
	// return &db.Connector{DB: db, Ctx: context.Background()}
}

func getDB(name string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(name), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
}
