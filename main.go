package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	file := "/tmp/teamux.log"
	logfile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Unable to setup syslog:", err.Error())
	}
	defer logfile.Close()
	logger := common.Logger{
		Infologger:    log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warninglogger: log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Errorlogger:   log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Fatallogger:   log.New(logfile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	p := tea.NewProgram(components.InitialModel(db, logger))

	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
