// Package panes defines the UI components to manage and interact with TMUX
// panes and saved layouts.
package panes

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/list"
)

type (
	PaneBrowserModel struct {
		list     list.Model
		selected string
		state    common.State
		db       *sql.DB
		logger   common.Logger
	}
	paneItem struct {
		title string
		desc  string
	}
)

func (pi paneItem) FilterValue() string {
	return ""
}
