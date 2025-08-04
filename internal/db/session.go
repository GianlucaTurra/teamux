package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type session struct {
	db   *sql.DB
	Id   int
	File string
	Name string
}

func (s session) NewSession(file string, name string, db *sql.DB) session {
	return session{
		db:   db,
		Id:   0,
		File: file,
		Name: name,
	}
}

func (s session) Save() error {
	var query string
	if s.Id == 0 {
		query = "INSERT INTO session_files (file, name) VALUES (?, ?)"
	} else {
		query = "UPDATE session_files SET file = ?, name = ? WHERE id = ?"
	}
	if _, err := s.db.Exec(query, s.File, s.Name, s.Id); err != nil {
		return err
	}
	return nil
}

func (s session) Read(id int) error {
	query := "SELECT id, file, name FROM session_files WHERE id = ?"
	row := s.db.QueryRow(query, id)
	if err := row.Scan(&s.Id, &s.File, &s.Name); err != nil {
		return err
	}
	return nil
}

func ReadAllSessions(db *sql.DB) ([]session, error) {
	query := "SELECT id, file, name FROM session_files"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []session
	for rows.Next() {
		var s session
		if err := rows.Scan(&s.Id, &s.File, &s.Name); err != nil {
			return nil, err
		}
		s.db = db // Assign the db to the session
		sessions = append(sessions, s)
	}
	return sessions, nil
}
