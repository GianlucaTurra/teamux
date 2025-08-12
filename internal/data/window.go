package data

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Window struct {
	db               *sql.DB
	ID               int
	Name             string
	WorkingDirectory string
}

func NewWindow(name string, workingDirectory string, db *sql.DB) Window {
	return Window{
		db:               db,
		ID:               0,
		Name:             name,
		WorkingDirectory: workingDirectory,
	}
}

func (w Window) Save() error {
	var query string
	if w.ID == 0 {
		query = "INSERT INTO windows (name, working_directory) VALUES (?, ?)"
	} else {
		query = "UPDATE windows SET name = ?, working_directory = ? WHERE id = ?"
	}
	if _, err := w.db.Exec(query, w.Name, w.WorkingDirectory, w.ID); err != nil {
		return err
	}
	return nil
}

func ReadAllWindows(db *sql.DB) ([]Window, error) {
	query := "SELECT id, name, working_directory FROM windows"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var windows []Window
	for rows.Next() {
		var w Window
		if err := rows.Scan(&w.ID, &w.Name, &w.WorkingDirectory); err != nil {
			return nil, err
		}
		if strings.TrimSpace(w.WorkingDirectory) == "" {
			if home, err := os.UserHomeDir(); err != nil {
				return nil, err
			} else {
				w.WorkingDirectory = home
			}
		}
		w.db = db // Assign the db to the window
		windows = append(windows, w)
	}
	return windows, nil
}

func ReadWindowByID(db *sql.DB, id int) (*Window, error) {
	query := "SELECT id, name, working_directory FROM windows WHERE id = ?"
	row := db.QueryRow(query, id)
	var w Window
	if err := row.Scan(&w.ID, &w.Name, &w.WorkingDirectory); err != nil {
		return nil, err
	}
	if strings.TrimSpace(w.WorkingDirectory) == "" {
		if home, err := os.UserHomeDir(); err != nil {
			return nil, err
		} else {
			w.WorkingDirectory = home
		}
	}
	w.db = db // Assign the db to the window
	return &w, nil
}

func (w Window) Delete() error {
	query := "DELETE FROM windows WHERE id = ?"
	if _, err := w.db.Exec(query, w.ID); err != nil {
		return err
	}
	return nil
}

func (w Window) Open() error {
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("tmux neww -d -n \"%s\" -c %s", w.Name, w.WorkingDirectory),
	)
	return cmd.Run()
}

func (w Window) OpenWithTarget(target string) error {
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("tmux neww -t %s -d -n \"%s\" -c %s", target, w.Name, w.WorkingDirectory),
	)
	return cmd.Run()
}
