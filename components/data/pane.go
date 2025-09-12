package data

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Pane struct {
	db               *sql.DB
	ID               int
	Name             string
	WorkingDirectory string
	splitDirection   int
	SplitRatio       int
}

const (
	vertical = iota
	horizontal
)

func NewVerticalPane(
	name string,
	workingDirectory string,
	splitRatio int,
	db *sql.DB,
) Pane {
	return newPane(name, workingDirectory, vertical, splitRatio, db)
}

func NewHorizontalPane(
	name string,
	workingDirectory string,
	splitRatio int,
	db *sql.DB,
) Pane {
	return newPane(name, workingDirectory, horizontal, splitRatio, db)
}

func newPane(
	name string,
	workingDirectory string,
	splitDirection int,
	splitRatio int,
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
		query = insertPane
	} else {
		query = updatePane
	}
	if _, err := p.db.Exec(
		query,
		p.Name,
		p.WorkingDirectory,
		p.splitDirection,
		p.SplitRatio,
		p.ID,
	); err != nil {
		return err
	}
	return nil
}

func (p Pane) Delete() error {
	if _, err := p.db.Exec(deletePaneByID, p.ID); err != nil {
		return err
	}
	return nil
}

func (p *Pane) SetHorizontal() { p.splitDirection = horizontal }

func (p *Pane) SetVertical() { p.splitDirection = vertical }

func (p Pane) Open(target *string) error {
	var tmuxCommand string
	if target != nil {
		tmuxCommand = fmt.Sprintf(
			"tmux split-window -t \"%s\" -l %s -c \"%s\"",
			*target,
			strconv.Itoa(p.SplitRatio)+"%",
			p.WorkingDirectory,
		)
	} else {
		tmuxCommand = fmt.Sprintf(
			"tmux split-window -l %s -c \"%s\"",
			strconv.Itoa(p.SplitRatio)+"%",
			p.WorkingDirectory,
		)
	}
	if p.IsHorizontal() {
		tmuxCommand += " -h"
	} else {
		tmuxCommand += " -v"
	}
	cmd := exec.Command("sh", "-c", tmuxCommand)
	return cmd.Run()
}

func GetAllPanes(db *sql.DB) ([]Pane, error) {
	rows, err := db.Query(selectAllPanes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var panes []Pane
	for rows.Next() {
		var p Pane
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.WorkingDirectory,
			&p.splitDirection,
			&p.SplitRatio,
		); err != nil {
			return nil, err
		}
		p.db = db
		// TODO: handle this error
		p.WorkingDirectory, _ = GetPWDorPlaceholder(p.WorkingDirectory)
		panes = append(panes, p)
	}
	return panes, nil
}

func GetPaneByID(db *sql.DB, id int) (*Pane, error) {
	row := db.QueryRow(selectPaneByID, id)
	var p Pane
	if err := row.Scan(
		&p.ID,
		&p.Name,
		&p.WorkingDirectory,
		&p.splitDirection,
		&p.SplitRatio,
	); err != nil {
		return nil, err
	}
	p.db = db
	// TODO: handle this error
	p.WorkingDirectory, _ = GetPWDorPlaceholder(p.WorkingDirectory)
	return &p, nil
}
