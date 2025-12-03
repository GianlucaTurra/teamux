// Package panes defines the UI components to manage and interact with TMUX
// panes and saved layouts.
package panes

import (
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type (
	PaneBrowserModel struct {
		list      list.Model
		data      map[string]data.Pane
		selected  string
		state     common.State
		connector data.Connector
		logger    common.Logger
	}
	paneItem struct {
		title string
	}
)

func (pi paneItem) FilterValue() string {
	return ""
}

// loadPaneData Load the panes from db into the current model
func loadPaneData(db *gorm.DB, logger common.Logger) (map[string]data.Pane, []list.Item) {
	layouts := []list.Item{}
	paneData := make(map[string]data.Pane)
	panes, err := data.ReadAllPanes(db)
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

func NewPaneBrowserModel(connector data.Connector, logger common.Logger) PaneBrowserModel {
	data, layouts := loadPaneData(connector.DB, logger)
	l := list.New(layouts, paneDelegate{}, 100, 10)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	return PaneBrowserModel{
		list:      l,
		data:      data,
		state:     common.Browsing,
		logger:    logger,
		connector: connector,
	}
}

func (m PaneBrowserModel) Init() tea.Cmd {
	return nil
}

func (m PaneBrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case common.ReloadMsg:
		return NewPaneBrowserModel(m.connector, m.logger), nil
	case common.DeleteMsg:
		return m.deleteSelected()
	case common.OpenMsg:
		return m.openSelected()
	case tea.KeyMsg:
		if m.state == common.Deleting {
			return m.confirmDeletion(msg)
		}
		switch msg.String() {
		case "enter", " ":
			return m.open()
		case "q", "esc", "ctrl+c":
			m.state = common.Quitting
			return m, common.Quit
		case "d":
			return m.delete()
		case "e":
			return m.edit()
		case "n":
			return m, func() tea.Msg { return common.NewPaneMsg{} }
		}
	}
	m.list, cmd = m.list.Update(msg)
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

// edit() Open the current pane into the editor
func (m PaneBrowserModel) edit() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(paneItem)
	if ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return common.EditP(m.data[m.selected]) }
}

// delete Mark the current pane for deletion signaling for the confirmation
// message to appear
func (m PaneBrowserModel) delete() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(paneItem)
	if ok {
		m.selected = i.title
	}
	m.state = common.Deleting
	return m, nil
}

// open Signals to open the selected pane into the current tmux window
func (m PaneBrowserModel) open() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(paneItem)
	if ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return common.OpenMsg{} }
}

// confirmDeletion Delete from the db the pane marked for deletion
func (m PaneBrowserModel) confirmDeletion(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.state = common.Browsing
		return m, common.Delete
	default:
		m.state = common.Browsing
		return m, nil
	}
}

// openSelected Run the tmux command to open the selected pane
func (m PaneBrowserModel) openSelected() (PaneBrowserModel, tea.Cmd) {
	p := m.data[m.selected]
	if err := p.Open(); err != nil {
		e := fmt.Errorf("failed to open pane %s: %v", m.selected, err)
		m.logger.Errorlogger.Println(e.Error())
		return m, func() tea.Msg { return common.OutputMsg{Err: e, Severity: common.Error} }
	}
	return m, func() tea.Msg { return common.ReloadMsg{} }
}

// deleteSelected Delete the selected pane from the db
func (m PaneBrowserModel) deleteSelected() (PaneBrowserModel, tea.Cmd) {
	p := m.data[m.selected]
	if _, err := p.Delete(m.connector); err != nil {
		e := fmt.Errorf("failed to delete pane %s: %v", m.selected, err)
		m.logger.Errorlogger.Println(e.Error())
		return m, func() tea.Msg { return common.OutputMsg{Err: e, Severity: common.Error} }
	}
	return m, func() tea.Msg { return common.ReloadMsg{} }
}
