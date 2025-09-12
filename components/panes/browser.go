// Package panes defines the UI components to manage and interact with TMUX
// panes and saved layouts.
package panes

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	PaneBrowserModel struct {
		list     list.Model
		data     map[string]data.Pane
		selected string
		state    common.State
		db       *sql.DB
		logger   common.Logger
	}
	paneItem struct {
		title string
	}
)

func (pi paneItem) FilterValue() string {
	return ""
}

func loadPaneData(db *sql.DB, logger common.Logger) (map[string]data.Pane, []list.Item) {
	layouts := []list.Item{}
	paneData := make(map[string]data.Pane)
	panes, err := data.GetAllPanes(db)
	if err != nil {
		logger.Fatallogger.Fatalf("Failed to read panes: %v", err)
		return paneData, layouts
	}
	for _, p := range panes {
		paneData[p.Name] = p
		layouts = append(layouts, paneItem{title: p.Name})
	}
	return paneData, layouts
}

func NewPaneBrowserModel(db *sql.DB, logger common.Logger) PaneBrowserModel {
	data, layouts := loadPaneData(db, logger)
	l := list.New(layouts, paneDelegate{}, 100, 10)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	return PaneBrowserModel{
		list:   l,
		data:   data,
		state:  common.Browsing,
		logger: logger,
		db:     db,
	}
}

func (m PaneBrowserModel) Init() tea.Cmd {
	return nil
}

func (m PaneBrowserModel) Update(msg tea.Msg) (PaneBrowserModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.state = common.Quitting
			return m, tea.Quit
		}
	}
	m.list, cmd = m.list.Update(msg)
	m.selected = m.list.SelectedItem().(paneItem).title
	return m, cmd
}

func (m PaneBrowserModel) View() string {
	switch m.state {
	case common.Quitting:
		return ""
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.list.View(),
	)
}
