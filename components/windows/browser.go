// Package windows declares the UI models to show and interact with saved
// TMUX window layouts
package windows

import (
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/GianlucaTurra/teamux/tmux"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	windowItem struct {
		title string
		desc  string
	}
	WindowBrowserModel struct {
		list      list.Model
		selected  string
		state     common.State
		data      map[string]data.Window
		connector data.Connector
		logger    common.Logger
	}
)

func (s windowItem) FilterValue() string { return "" }

func NewWindowBrowserModel(connector data.Connector, logger common.Logger) WindowBrowserModel {
	data, layouts := loadWindowData(connector, logger)
	l := list.New(layouts, WindowDelegate{}, 100, 10)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	return WindowBrowserModel{
		list:      l,
		data:      data,
		state:     common.Browsing,
		logger:    logger,
		connector: connector,
	}
}

func loadWindowData(connector data.Connector, logger common.Logger) (map[string]data.Window, []list.Item) {
	layouts := []list.Item{}
	windowData := make(map[string]data.Window)
	windows, err := data.ReadAllWindows(connector.DB)
	if err != nil {
		logger.Fatallogger.Fatalf("Failed to read windows: %v", err)
		return windowData, layouts
	}
	for _, w := range windows {
		layouts = append(layouts, windowItem{title: w.Name})
		windowData[w.Name] = w
	}
	return windowData, layouts
}

func (m WindowBrowserModel) View() string {
	switch m.state {
	case common.Quitting:
		return "Bye, have a nice day!"
	case common.Deleting:
		return fmt.Sprintf("You are about to delete %s, press y to confirm", m.selected)
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.list.View(),
	)
}

func (m WindowBrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case common.OpenMsg:
		return m, func() tea.Msg { return openSelected(m.logger, m.data[m.selected]) }
	case common.DeleteMsg:
		return m, func() tea.Msg { return deleteSelected(m.logger, m.data[m.selected], m.connector) }
	case common.KillMsg:
		return m, func() tea.Msg { return killSelected(m.logger, m.data[m.selected]) }
	case common.ReloadMsg:
		return NewWindowBrowserModel(m.connector, m.logger), nil
	case common.UpDownMsg:
		return m.selectUpDown()
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
		case "ctrl+c", "q", "esc":
			m.state = common.Quitting
			return m, common.Quit
		case "enter", " ":
			return m.open()
		case "s":
			return m.switchToSelected()
		case "d":
			return m.delete()
		// TODO: for consistency with the tmux shortcuts use x instead
		case "K":
			return m.kill()
		case "e":
			return m.editSelected()
		case "a":
			return m.addPanesToSelected()
		case "n":
			return m, func() tea.Msg { return common.NewWindowMsg{} }
		case "j", "k", "up", "down":
			cmds = append(cmds, common.UpDown)
		case "?":
			return m, func() tea.Msg { return common.ShowFullHelpMsg{Component: common.WindowBrowser} }
		}
	}
	// handle sub-models updates
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m WindowBrowserModel) Init() tea.Cmd {
	return nil
}

func (m WindowBrowserModel) addPanesToSelected() (tea.Model, tea.Cmd) {
	if i, ok := m.list.SelectedItem().(windowItem); ok {
		window := m.data[i.title]
		return m, func() tea.Msg { return common.AssociatePanesMsg{Window: window} }
	}
	return nil, nil
}

func (m WindowBrowserModel) editSelected() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(windowItem)
	if ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return common.EditW(m.data[m.selected]) }
}

func (m WindowBrowserModel) kill() (tea.Model, tea.Cmd) {
	if i, ok := m.list.SelectedItem().(windowItem); ok {
		m.selected = i.title
	}
	return m, common.Kill
}

func (m WindowBrowserModel) delete() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(windowItem)
	if ok {
		m.selected = i.title
	}
	m.state = common.Deleting
	return m, nil
}

func (m WindowBrowserModel) switchToSelected() (tea.Model, tea.Cmd) {
	if i, ok := m.list.SelectedItem().(windowItem); ok {
		m.selected = i.title
	}
	return m, common.Switch
}

func (m WindowBrowserModel) open() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(windowItem)
	if ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return common.OpenMsg{} }
}

func (m WindowBrowserModel) selectUpDown() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(windowItem)
	if ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return common.NewWFocus{Window: m.data[m.selected]} }
}

func openSelected(logger common.Logger, w data.Window) tea.Cmd {
	if err := w.Open(); err != nil {
		logger.Errorlogger.Printf("Error opening window %s: %v", w.Name, err)
		var severity common.Severity
		switch err.(type) {
		case tmux.Warning:
			severity = common.Warning
		default:
			severity = common.Error
		}
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: severity} }
	}
	return nil
}

func deleteSelected(logger common.Logger, w data.Window, connector data.Connector) tea.Cmd {
	if _, err := w.Delete(connector); err != nil {
		logger.Errorlogger.Printf("Error deleting window %s: %v", w.Name, err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	return func() tea.Msg { return common.ReloadMsg{} }
}

func killSelected(logger common.Logger, w data.Window) tea.Cmd {
	if err := w.Kill(); err != nil {
		logger.Errorlogger.Printf("Error killing window %s: %v", w.Name, err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	return func() tea.Msg { return common.ReloadMsg{} }
}
