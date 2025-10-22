// Package panes defines the UI components to manage and interact with TMUX
// panes and saved layouts.
package panes

import (
	"database/sql"
	"fmt"

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

func NewPaneBrowserModel(db *sql.DB, logger common.Logger) common.TeamuxModel {
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

func (m PaneBrowserModel) Update(msg tea.Msg) (common.TeamuxModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case common.ReloadMsg:
		return NewPaneBrowserModel(m.db, m.logger), nil
	case common.DeleteMsg:
		return m.deleteSelected()
	case common.OpenMsg:
		return m.openSelected()
	case tea.KeyMsg:
		if m.state == common.Deleting {
			switch msg.String() {
			case "y":
				m.state = common.Browsing
				return m, common.Delete
			default:
				m.state = common.Browsing
				return m, nil
			}
		}
		switch msg.String() {
		case "enter", " ":
			i, ok := m.list.SelectedItem().(paneItem)
			if ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return common.OpenMsg{} }
		case "q", "esc", "ctrl+c":
			m.state = common.Quitting
			return m, tea.Quit
		case "d":
			i, ok := m.list.SelectedItem().(paneItem)
			if ok {
				m.selected = i.title
			}
			m.state = common.Deleting
			return m, nil
		case "e":
			i, ok := m.list.SelectedItem().(paneItem)
			if ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return common.EditP(m.data[m.selected]) }
		case "n":
			return m, func() tea.Msg { return common.NewPaneMsg{} }
		}
	}
	m.list, cmd = m.list.Update(msg)
	m.selected = m.list.SelectedItem().(paneItem).title
	return m, cmd
}

func (m PaneBrowserModel) View() string {
	switch m.state {
	case common.Deleting:
		return fmt.Sprintf("You are about to delete %s, press y to confirm", m.selected)
	case common.Quitting:
		return ""
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.list.View(),
	)
}

func (m PaneBrowserModel) openSelected() (PaneBrowserModel, tea.Cmd) {
	p := m.data[m.selected]
	if err := p.Open(nil); err != nil {
		m.logger.Errorlogger.Printf("Failed to open pane %s: %v", m.selected, err)
		return m, nil
	}
	return m, func() tea.Msg { return common.ReloadMsg{} }
}

func (m PaneBrowserModel) deleteSelected() (PaneBrowserModel, tea.Cmd) {
	p := m.data[m.selected]
	if err := p.Delete(); err != nil {
		m.logger.Errorlogger.Printf("Failed to delete pane %s: %v", m.selected, err)
		// TODO: actually do something
		return m, nil
	}
	return m, func() tea.Msg { return common.ReloadMsg{} }
}

func (m PaneBrowserModel) GetDB() *sql.DB {
	return m.db
}

func (m PaneBrowserModel) GetLogger() common.Logger {
	return m.logger
}
