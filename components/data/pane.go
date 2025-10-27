package data

import (
	"fmt"
	"os/exec"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type Pane struct {
	gorm.Model
	Name             string
	WorkingDirectory string
	splitDirection   int
	SplitRatio       int
}

const (
	vertical = iota
	horizontal
)

func ReadAllPanes(db *gorm.DB) ([]Pane, error) {
	var panes []Pane
	err := db.Model(&Pane{}).Find(panes).Error
	return panes, err
}

func CreateVerticalPane(
	name string,
	workingDirectory string,
	splitRatio int,
	connector Connector,
) (int64, error) {
	return createPane(name, workingDirectory, vertical, splitRatio, connector)
}

func CreateHorizontalPane(
	name string,
	workingDirectory string,
	splitRatio int,
	connector Connector,
) (int64, error) {
	return createPane(name, workingDirectory, horizontal, splitRatio, connector)
}

func createPane(
	name string,
	workingDirectory string,
	splitDirection int,
	splitRatio int,
	connector Connector,
) (int64, error) {
	pane := Pane{
		Name:             name,
		WorkingDirectory: workingDirectory,
		splitDirection:   splitDirection,
		SplitRatio:       splitRatio,
	}
	result := gorm.WithResult()
	err := gorm.G[Pane](connector.DB, result).Create(connector.Ctx, &pane)
	return result.RowsAffected, err
}

func (p Pane) IsVertical() bool { return p.splitDirection == vertical }

func (p Pane) IsHorizontal() bool { return p.splitDirection == horizontal }

func (p Pane) Save(connector Connector) (int, error) {
	return gorm.G[Pane](connector.DB).Updates(connector.Ctx, p)
}

func (p Pane) Delete(connector Connector) (int, error) {
	return gorm.G[Pane](connector.DB).Where("id = ?", p.ID).Delete(connector.Ctx)
}

func (p *Pane) SetHorizontal() { p.splitDirection = horizontal }

func (p *Pane) SetVertical() { p.splitDirection = vertical }

// TODO: move to a proper package

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
