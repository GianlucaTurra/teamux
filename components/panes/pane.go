package panes

import (
	"github.com/GianlucaTurra/teamux/database"
	"github.com/GianlucaTurra/teamux/tmux"
	"gorm.io/gorm"
)

type Pane struct {
	gorm.Model
	Name             string
	WorkingDirectory string
	SplitDirection   string `gorm:"default:v;"`
	SplitRatio       int
	Target           string
	ShellCmd         string
}

const (
	Vertical   = "v"
	Horizontal = "h"
)

func ReadAllPanes(db *gorm.DB) ([]Pane, error) {
	var panes []Pane
	err := db.Model(&Pane{}).Find(&panes).Error
	return panes, err
}

func CreateVerticalPane(
	name string,
	workingDirectory string,
	splitRatio int,
	connector database.Connector,
	target string,
	shellCmd string,
) (int64, error) {
	return createPane(name, workingDirectory, Vertical, splitRatio, connector, target, shellCmd)
}

func CreateHorizontalPane(
	name string,
	workingDirectory string,
	splitRatio int,
	connector database.Connector,
	target string,
	shellCmd string,
) (int64, error) {
	return createPane(name, workingDirectory, Horizontal, splitRatio, connector, target, shellCmd)
}

func createPane(
	name string,
	workingDirectory string,
	splitDirection string,
	splitRatio int,
	connector database.Connector,
	target string,
	shellCmd string,
) (int64, error) {
	pane := Pane{
		Name:             name,
		WorkingDirectory: workingDirectory,
		SplitDirection:   splitDirection,
		SplitRatio:       splitRatio,
		Target:           target,
		ShellCmd:         shellCmd,
	}
	result := gorm.WithResult()
	err := gorm.G[Pane](connector.DB, result).Create(connector.Ctx, &pane)
	return result.RowsAffected, err
}

func (p Pane) Save(connector database.Connector) (int, error) {
	return gorm.G[Pane](connector.DB).Updates(connector.Ctx, p)
}

func (p Pane) Delete(connector database.Connector) (int, error) {
	return gorm.G[Pane](connector.DB).Where("id = ?", p.ID).Delete(connector.Ctx)
}

func (p Pane) Open() error {
	return tmux.SplitWindow(p.SplitRatio, p.WorkingDirectory, p.SplitDirection, p.ShellCmd)
}

func (p Pane) OpenWithTarget(target string) error {
	return tmux.SplitWindowWithTargetWindow(target, p.SplitRatio, p.WorkingDirectory, p.SplitDirection, p.ShellCmd)
}
