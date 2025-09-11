package data

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Pane struct {
	db               *sql.DB
	ID               int
	Name             string
	WorkingDirectory string
	splitDirection   int
	SplitRatio       float32
}

const (
	vertical = iota
	horizontal
)

func NewVerticalPane(
	name string,
	workingDirectory string,
	splitRatio float32,
	db *sql.DB,
) Pane {
	return newPane(name, workingDirectory, vertical, splitRatio, db)
}

func NewHorizontalPane(
	name string,
	workingDirectory string,
	splitRatio float32,
	db *sql.DB,
) Pane {
	return newPane(name, workingDirectory, horizontal, splitRatio, db)
}

func newPane(
	name string,
	workingDirectory string,
	splitDirection int,
	splitRatio float32,
	db *sql.DB,
) Pane {
	return Pane{
		db:               db,
		ID:               0,
		Name:             name,
		WorkingDirectory: workingDirectory,
		splitDirection:   splitDirection,
		SplitRatio:       splitRatio,
	}
}

func (p Pane) IsVertical() bool { return p.splitDirection == vertical }

func (p Pane) IsHorizontal() bool { return p.splitDirection == horizontal }

func (p Pane) Save() error {
	var query string
	if p.ID == 0 {
		query = `
			INSERT INTO 
			panes (name, working_directory, split_direction, split_ratio)
			VALUES (?, ?, ?, ?)
		`
	} else {
		query = `
			UPDATE panes 
			SET name = ?, working_directory = ?, split_direction = ?, split_ratio = ? 
			WHERE id = ?
		`
	}
	if _, err := p.db.Exec(query, p.Name, p.WorkingDirectory, p.splitDirection, p.SplitRatio, p.ID); err != nil {
		return err
	}
	return nil
}
