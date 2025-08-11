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
)

type (
	windowItem struct {
		title string
		desc  string
	}
	windowBrowserModel struct {
		list        list.Model
		selected    string
		windowState State
		data        map[string]data.Window
		db          *sql.DB
		logger      internal.Logger
		// help     windowBrowserHelpModel
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
		list:        l,
		data:        data,
		windowState: browsing,
		logger:      logger,
		db:          db,
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
