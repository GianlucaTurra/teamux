package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SessionsInfo struct {
	file string
	name string
}

func ReadSeassions() []SessionsInfo {
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
	var sessions []SessionsInfo
	for rows.Next() {
		var session SessionsInfo
		err = rows.Scan(&session)
		if err != nil {
			log.Fatal(err)
		}
		sessions = append(sessions, session)
	}
	return sessions
}
