// Package data declares the data structures to map db entities
package data

import (
	"database/sql"
	"fmt"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

type Session struct {
	db               *sql.DB
	ID               int
	Name             string
	WorkingDirectory string
	Windows          []Window
}

func NewSession(name string, workingDirectory string, db *sql.DB) Session {
	return Session{
		db:               db,
		ID:               0,
		Name:             name,
		WorkingDirectory: workingDirectory,
	}
}

// Save Create a new record on the database or update it if it already exists
func (s Session) Save() error {
	var query string
	if s.ID == 0 {
		query = insertSessions
	} else {
		query = updateSession
	}
	if _, err := s.db.Exec(query, s.Name, s.WorkingDirectory, s.ID); err != nil {
		return err
	}
	return nil
}

// ReadAllSessions from the database. If a session is missing the working
// directory it is set to the user's home directory.
func ReadAllSessions(db *sql.DB) ([]Session, error) {
	rows, err := db.Query(selectAllSessions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		s.db = db
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func ReadSessionByID(db *sql.DB, id int) (*Session, error) {
	row := db.QueryRow(selectSessionByID, id)
	var s Session
	if err := row.Scan(&s.ID, &s.Name, &s.WorkingDirectory); err != nil {
		return nil, err
	}
	s.db = db
	return &s, nil
}

// Delete removes a session from the database by its ID.
func (s Session) Delete() error {
	if _, err := s.db.Exec(deleteSessionByID, s.ID); err != nil {
		return err
	}
	return nil
}

// Open translates the session object to a tmux command to open a new session.
func (s Session) Open() error {
	newSessionCmd := fmt.Sprintf("tmux new-session -d -s \"%s\" -c %s", s.Name, s.WorkingDirectory)
	cmd := exec.Command("sh", "-c", newSessionCmd)
	return cmd.Run()
}

func (s Session) IsOpen() bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux has-session -t \"%s\"", s.Name))
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (s Session) Close() error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux kill-session -t \"%s\"", s.Name))
	return cmd.Run()
}

func (s Session) Switch() error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux switch -t \"%s\"", s.Name))
	return cmd.Run()
}

// GetAllWindows reads all Window.ID from the association table based on the
// current Session.ID
func (s *Session) GetAllWindows() error {
	rows, err := s.db.Query(selectAllSessionWindows, s.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var window Window
		if err := rows.Scan(&window.ID, &window.Name, &window.WorkingDirectory); err != nil {
			return err
		}
		window.db = s.db
		s.Windows = append(s.Windows, window)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (s *Session) GetPWD() error {
	row, err := s.db.Query(selectSessionWorkingDirectory, s.ID)
	if err != nil {
		return err
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&s.WorkingDirectory); err != nil {
			return err
		}
	}
	checkedPath, err := getPWD(s.WorkingDirectory)
	if err != nil {
		return err
	}
	s.WorkingDirectory = checkedPath
	return nil
}

func GetFirstSession(db *sql.DB) (Session, error) {
	row, err := db.Query(selectFirstSession)
	if err != nil {
		return Session{}, err
	}
	defer row.Close()
	var session Session
	for row.Next() {
		if err := row.Scan(&session.ID, &session.Name, &session.WorkingDirectory); err != nil {
			return Session{}, err
		}
		session.db = db
	}
	return session, nil
}
