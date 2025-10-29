package data

import (
	"github.com/GianlucaTurra/teamux/tmux"
	"gorm.io/gorm"
)

type Window struct {
	gorm.Model
	Name             string
	WorkingDirectory string
	Panes            []Pane `gorm:"many2many:window_panes"`
}

func CreateWindow(name string, workingDirectory string, connector Connector) (int64, error) {
	window := Window{Name: name, WorkingDirectory: workingDirectory}
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
	err := db.Model(&Window{}).Preload("Panes").Find(windows).Error
	return windows, err
}

// TODO: remove these methods

func (w Window) Open() error {
	return tmux.CreateWindow(w.Name, w.WorkingDirectory)
}

func (w Window) Kill() error {
	return tmux.KillWindow(w.Name)
}

func (w Window) OpenWithTarget(target string) error {
	return tmux.CreateWindowWithTarget(w.Name, w.WorkingDirectory, target)
}
