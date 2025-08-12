package components

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/GianlucaTurra/teamux/internal"
	"github.com/GianlucaTurra/teamux/internal/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	windowItem struct {
		title string
		desc  string
	}
	windowBrowserModel struct {
		list     list.Model
		selected string
		state    State
		data     map[string]data.Window
		db       *sql.DB
		logger   internal.Logger
		help     windowBrowserHelpModel
	}
	WindowDelegate struct{}
)

func (s windowItem) FilterValue() string { return "" }

func (d WindowDelegate) Height() int                             { return 1 }
func (d WindowDelegate) Spacing() int                            { return 0 }
func (d WindowDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d WindowDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(windowItem)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

func newWindowBrowserModel(db *sql.DB, logger internal.Logger) windowBrowserModel {
	data, layouts := loadWindowData(db, logger)
	l := list.New(layouts, WindowDelegate{}, 100, 10)
	l.Title = "Available window layouts"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return windowBrowserModel{
		list:   l,
		data:   data,
		state:  browsing,
		logger: logger,
		db:     db,
		help:   newWindowBrowserHelpModel(),
	}
}

func loadWindowData(db *sql.DB, logger internal.Logger) (map[string]data.Window, []list.Item) {
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

func (m windowBrowserModel) View() string {
	switch m.state {
	case quitting:
		return "Bye, have a nice day!"
	case deleting:
		return fmt.Sprintf("You are about to delete %s, press y to confirm", m.selected)
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.list.View(),
		m.help.View(),
	)
}

func (m windowBrowserModel) Update(msg tea.Msg) (windowBrowserModel, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.OpenMsg:
		return m.openSelected()
	case internal.DeleteMsg:
		return m.deleteSelected()
	case tea.KeyMsg:
		if m.state == deleting {
			switch msg.String() {
			case "y":
				m.state = browsing
				return m, internal.Delete
			default:
				m.state = browsing
				return m, nil
			}
		}
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.state = quitting
			return m, tea.Quit
		case "enter", " ":
			i, ok := m.list.SelectedItem().(windowItem)
			if ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return internal.OpenMsg{} }
		case "s":
			if i, ok := m.list.SelectedItem().(windowItem); ok {
				m.selected = i.title
			}
			return m, internal.Switch
		case "d":
			i, ok := m.list.SelectedItem().(windowItem)
			if ok {
				m.selected = i.title
			}
			m.state = deleting
			return m, nil
		case "K":
			if i, ok := m.list.SelectedItem().(windowItem); ok {
				m.selected = i.title
			}
			return m, internal.Kill
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

func (m windowBrowserModel) Init() tea.Cmd {
	return nil
}

// openSelected Opens the selected window.
func (m windowBrowserModel) openSelected() (windowBrowserModel, tea.Cmd) {
	w := m.data[m.selected]
	if err := w.Open(); err != nil {
		m.logger.Errorlogger.Printf("Error opening window %s: %v", w.Name, err)
		return m, func() tea.Msg { return internal.TmuxErr{} }
	}
	return m, func() tea.Msg { return nil }
}

// deleteSelected delete the window from the db
func (m windowBrowserModel) deleteSelected() (windowBrowserModel, tea.Cmd) {
	w := m.data[m.selected]
	if err := w.Delete(); err != nil {
		m.logger.Errorlogger.Printf("Error deleting window %s: %v", m.selected, err)
		return m, func() tea.Msg { return internal.TmuxErr{} }
	}
	return m, func() tea.Msg { return internal.ReloadMsg{} }
}
