package windows

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
	windowItem struct {
		title string
		desc  string
	}
	WindowBrowserModel struct {
		list     list.Model
		selected string
		state    common.State
		data     map[string]data.Window
		db       *sql.DB
		logger   common.Logger
		help     windowBrowserHelpModel
	}
)

func (s windowItem) FilterValue() string { return "" }

func NewWindowBrowserModel(db *sql.DB, logger common.Logger) WindowBrowserModel {
	data, layouts := loadWindowData(db, logger)
	l := list.New(layouts, WindowDelegate{}, 100, 10)
	l.SetShowTitle(false)
	l.Styles.Title = common.TitleStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	l.Styles.HelpStyle = common.HelpStyle
	return WindowBrowserModel{
		list:   l,
		data:   data,
		state:  common.Browsing,
		logger: logger,
		db:     db,
		help:   newWindowBrowserHelpModel(),
	}
}

func loadWindowData(db *sql.DB, logger common.Logger) (map[string]data.Window, []list.Item) {
	layouts := []list.Item{}
	windowData := make(map[string]data.Window)
	windows, err := data.ReadAllWindows(db)
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
		m.help.View(),
	)
}

func (m WindowBrowserModel) Update(msg tea.Msg) (WindowBrowserModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.OpenMsg:
		return m.openSelected()
	case common.DeleteMsg:
		return m.deleteSelected()
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
			i, ok := m.list.SelectedItem().(windowItem)
			if ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return common.OpenMsg{} }
		case "s":
			if i, ok := m.list.SelectedItem().(windowItem); ok {
				m.selected = i.title
			}
			return m, common.Switch
		case "d":
			i, ok := m.list.SelectedItem().(windowItem)
			if ok {
				m.selected = i.title
			}
			m.state = common.Deleting
			return m, nil
		case "K":
			if i, ok := m.list.SelectedItem().(windowItem); ok {
				m.selected = i.title
			}
			return m, common.Kill
		}
	}
	// handle sub-models updates
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	newHelp, cmd := m.help.Update(msg)
	m.help = newHelp
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m WindowBrowserModel) Init() tea.Cmd {
	return nil
}

// openSelected Opens the selected window.
func (m WindowBrowserModel) openSelected() (WindowBrowserModel, tea.Cmd) {
	w := m.data[m.selected]
	if err := w.Open(); err != nil {
		m.logger.Errorlogger.Printf("Error opening window %s: %v", w.Name, err)
		return m, func() tea.Msg { return common.TmuxErr{} }
	}
	return m, func() tea.Msg { return nil }
}

// deleteSelected delete the window from the db
func (m WindowBrowserModel) deleteSelected() (WindowBrowserModel, tea.Cmd) {
	w := m.data[m.selected]
	if err := w.Delete(); err != nil {
		m.logger.Errorlogger.Printf("Error deleting window %s: %v", m.selected, err)
		return m, func() tea.Msg { return common.TmuxErr{} }
	}
	return m, func() tea.Msg { return common.ReloadMsg{} }
}
