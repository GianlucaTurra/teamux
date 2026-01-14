package data

import (
	"fmt"

	"github.com/GianlucaTurra/teamux/tmux"
	"gorm.io/gorm"
)

type Window struct {
	gorm.Model
	Name             string
	WorkingDirectory string
	ShellCmd         string
	Panes            []Pane `gorm:"many2many:window_panes"`
}

func CreateWindow(name string, workingDirectory string, shellCmd string, connector Connector) (int64, error) {
	window := Window{Name: name, WorkingDirectory: workingDirectory, ShellCmd: shellCmd}
	result := gorm.WithResult()
	err := gorm.G[Window](connector.DB, result).Create(connector.Ctx, &window)
	return result.RowsAffected, err
}

func (w Window) Save(connector Connector) (int, error) {
	return gorm.G[Window](connector.DB).Updates(connector.Ctx, w)
}

func (w Window) Delete(connector Connector) (int, error) {
	return gorm.G[Window](connector.DB).Where("id = ?", w.ID).Delete(connector.Ctx)
}

func ReadAllWindows(db *gorm.DB) ([]Window, error) {
	var windows []Window
	err := db.Model(&Window{}).Preload("Panes").Find(&windows).Error
	return windows, err
}

func (w Window) Open() error {
	return w.openAndCascade(nil)
}

func (w Window) Kill() error {
	return tmux.KillWindow(w.Name)
}

func (w Window) OpenWithTarget(target string) error {
	return w.openAndCascade(&target)
}

func (w Window) openAndCascade(target *string) error {
	var err error
	err = tmux.NewWindow(w.Name, w.WorkingDirectory, w.ShellCmd, target)
	if err != nil {
		return err
	}
	var qualifiedTarget string
	if target == nil {
		qualifiedTarget = w.Name
	} else {
		qualifiedTarget = *target + ":" + w.Name
	}
	for _, pane := range w.Panes {
		err = pane.OpenWithTarget(qualifiedTarget)
		if err != nil {
			err = tmux.NewWarning(fmt.Sprintf("error opening child pane %s: %v", pane.Name, err))
		}
	}
	return err
}
