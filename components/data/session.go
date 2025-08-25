package data

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Session struct {
	db               *sql.DB
	ID               int
	Name             string
	WorkingDirectory string
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
		query = "INSERT INTO sessions (name, working_directory) VALUES (?, ?)"
	} else {
		query = "UPDATE sessions SET name = ?, working_directory = ? WHERE id = ?"
	}
	if _, err := s.db.Exec(query, s.Name, s.WorkingDirectory, s.ID); err != nil {
		return err
	}
	return nil
}

// ReadAllSessions from the database. If a session is missing the working
// directory it is set to the user's home directory.
func ReadAllSessions(db *sql.DB) ([]Session, error) {
	query := "SELECT id, name, working_directory FROM sessions"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.ID, &s.Name, &s.WorkingDirectory); err != nil {
			return nil, err
		}
		if strings.TrimSpace(s.WorkingDirectory) == "" {
			if home, err := os.UserHomeDir(); err != nil {
				return nil, err
			} else {
				s.WorkingDirectory = home
			}
		}
		s.db = db // Assign the db to the session
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func ReadSessionByID(db *sql.DB, id int) (*Session, error) {
	query := "SELECT id, name, working_directory FROM sessions WHERE id = ?"
	row := db.QueryRow(query, id)
	var s Session
	if err := row.Scan(&s.ID, &s.Name, &s.WorkingDirectory); err != nil {
		return nil, err
	}
	if strings.TrimSpace(s.WorkingDirectory) == "" {
		if home, err := os.UserHomeDir(); err != nil {
			return nil, err
		} else {
			s.WorkingDirectory = home
		}
	}
	s.db = db // Assign the db to the session
	return &s, nil
}

// Delete removes a session from the database by its ID.
func (s Session) Delete() error {
	query := "DELETE FROM sessions WHERE id = ?"
	if _, err := s.db.Exec(query, s.ID); err != nil {
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
func (s Session) GetAllWindows() ([]int, error) {
	var windowsIds []int
	query := "SELECT window_id FROM Session_Windows WHERE session_id = ?"
	rows, err := s.db.Query(query, s.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		windowsIds = append(windowsIds, id)
	}
	if err = rows.Err(); err != nil {
		return windowsIds, err
	}
	return windowsIds, nil
}

func CountTmuxSessions() string {
	cmd := exec.Command("sh", "-c", "tmux ls | wc -l")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error counting sessions %v", err)
		return "Error executing tmux ls"
	}
	return out.String()
}
