package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SessionInfo struct {
	File string
	Name string
}

func ReadSeassions() map[string]string {
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
	sessions := make(map[string]string)
	for rows.Next() {
		var file string
		var name string
		err = rows.Scan(&file, &name)
		if err != nil {
			log.Fatal(err)
		}
		sessions[name] = file
	}
	return sessions
}
