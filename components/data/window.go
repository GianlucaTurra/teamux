package data

import (
	"fmt"
	"os/exec"

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
	err := db.Model(&Window{}).Preload("Panes").Find(&windows).Error
	return windows, err
}

// TODO: move these methods to a proper package

func (w Window) Open() error {
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("tmux neww -d -n \"%s\" -c %s", w.Name, w.WorkingDirectory),
	)
	return cmd.Run()
}

func (w Window) Kill() error {
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("tmux kill-window -t \"%s\"", w.Name),
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
