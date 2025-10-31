package data

import (
	"github.com/GianlucaTurra/teamux/tmux"
	"gorm.io/gorm"
)

type Pane struct {
	gorm.Model
	Name             string
	WorkingDirectory string
	splitDirection   int
	SplitRatio       int
	Target           string
}

const (
	vertical = iota
	horizontal
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
	connector Connector,
	target string,
) (int64, error) {
	return createPane(name, workingDirectory, vertical, splitRatio, connector, target)
}

func CreateHorizontalPane(
	name string,
	workingDirectory string,
	splitRatio int,
	connector Connector,
	target string,
) (int64, error) {
	return createPane(name, workingDirectory, horizontal, splitRatio, connector, target)
}

func createPane(
	name string,
	workingDirectory string,
	splitDirection int,
	splitRatio int,
	connector Connector,
	target string,
) (int64, error) {
	pane := Pane{
		Name:             name,
		WorkingDirectory: workingDirectory,
		splitDirection:   splitDirection,
		SplitRatio:       splitRatio,
		Target:           target,
	}
	result := gorm.WithResult()
	err := gorm.G[Pane](connector.DB, result).Create(connector.Ctx, &pane)
	return result.RowsAffected, err
}

func (p Pane) Save(connector Connector) (int, error) {
	return gorm.G[Pane](connector.DB).Updates(connector.Ctx, p)
}

func (p Pane) Delete(connector Connector) (int, error) {
	return gorm.G[Pane](connector.DB).Where("id = ?", p.ID).Delete(connector.Ctx)
}

func (p Pane) IsVertical() bool { return p.splitDirection == vertical }

func (p Pane) IsHorizontal() bool { return p.splitDirection == horizontal }

func (p *Pane) SetHorizontal() { p.splitDirection = horizontal }

func (p *Pane) SetVertical() { p.splitDirection = vertical }

// TODO: remove these methods

func (p Pane) Open() error {
	return tmux.CreatePane(p.Target, p.SplitRatio, p.WorkingDirectory, p.IsHorizontal())
}
