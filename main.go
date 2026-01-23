package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components"
	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	file := "/tmp/teamux.log"
	setup()
	logfile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalln("Unable to setup syslog:", err.Error())
	}
	defer logfile.Close()
	teamuxLogger := common.Logger{
		Infologger:    log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warninglogger: log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Errorlogger:   log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Fatallogger:   log.New(logfile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	db, err := database.GetProdDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(
		&sessions.Session{},
		&windows.Window{},
		&panes.Pane{},
	); err != nil {
		teamuxLogger.Errorlogger.Printf("Error migrating tables: %v", err)
	}
	p := tea.NewProgram(components.InitialModel(database.Connector{DB: db, Ctx: context.Background()}, teamuxLogger))
	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

func setup() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error detecting user config directory: %v", err)
	}
	teamuxConfigDir := fmt.Sprintf("%s/%s", userConfigDir, "teamux")
	if _, err := os.Stat(teamuxConfigDir); err == nil {
		return
	}
	if err := os.Mkdir(teamuxConfigDir, 0o755); err != nil {
		log.Fatalf("Error creating config dir: %v", err)
	}
}
