package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/syslog"
	"os"

	"github.com/GianlucaTurra/teamux/internal/components"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logWriter, err := syslog.New(syslog.LOG_SYSLOG, "teamux")
	if err != nil {
		log.Fatalln("Unable to setup syslog:", err.Error())
	}
	sysLogger := log.New(logWriter, "", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	p := tea.NewProgram(components.InitialModel(db, sysLogger))
	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
