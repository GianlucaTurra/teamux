package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components"
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	file := "/tmp/teamux.log"
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
	db, err := common.GetProdDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(
		&data.Session{},
		&data.Window{},
		&data.Pane{},
	); err != nil {
		teamuxLogger.Errorlogger.Printf("Error migrating tables: %v", err)
	}
	p := tea.NewProgram(components.InitialModel(data.Connector{DB: db, Ctx: context.Background()}, teamuxLogger))
	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
