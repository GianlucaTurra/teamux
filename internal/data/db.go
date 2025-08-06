package data

import (
	"database/sql"
	"log"

	"github.com/GianlucaTurra/teamux/internal"
	_ "github.com/mattn/go-sqlite3"
)

type SessionInfo struct {
	File   string
	IsOpen bool
}

func ReadSeassions() map[string]SessionInfo {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "select file, name from session_files"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	sessions := make(map[string]SessionInfo)
	for rows.Next() {
		var file string
		var name string
		err = rows.Scan(&file, &name)
		if err != nil {
			log.Fatal(err)
		}
		isOpen := internal.IsTmuxSessionOpen(name)
		sessions[name] = SessionInfo{file, isOpen}
	}
	return sessions
}
