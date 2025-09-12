package windows

import (
	"database/sql"
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WindowDetailModel struct {
	window data.Window
	db     *sql.DB
	logger common.Logger
}

func NewWindowDetailModel(db *sql.DB, logger common.Logger, window *data.Window) WindowDetailModel {
	if window == nil {
		firstWindow, err := data.GetFirstWindow(db)
		if err != nil {
			logger.Errorlogger.Printf("Error loading first window, falling back to default one.\n %v", err)
		}
		window = &firstWindow
	}
	if err := window.GetAllPanes(); err != nil {
		logger.Errorlogger.Printf("Error loading panes for first window.\n %v", err)
	}
	return WindowDetailModel{*window, db, logger}
}

func (m WindowDetailModel) View() string {
	var items []string
	title := renderTreeItem("Window Details", "", 0, false)
	items = append(items, common.TitleStyle.Foreground(lipgloss.Color("2")).Render(title))
	items = append(items, renderTreeItem(m.window.Name, m.window.WorkingDirectory, 0, false))
	for j, p := range m.window.Panes {
		items = append(
			items,
			renderTreeItem(p.Name, p.WorkingDirectory, 1, j == len(m.window.Panes)-1),
		)
	}
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func renderTreeItem(name string, pwd string, level int, isLast bool) string {
	var treeSymbol string
	if isLast {
		treeSymbol = "└──"
	} else {
		treeSymbol = "├──"
	}
	var prefix string
	switch level {
	case 0:
		prefix = ""
	case 1:
		prefix = treeSymbol + " "
	}
	nameCol := fmt.Sprintf("%-*s", 16, prefix+name)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		common.ItemStyle.Render(nameCol),
		common.ItemStyle.Italic(true).Render(pwd),
	)
}

func (m WindowDetailModel) Init() tea.Cmd { return nil }

func (m WindowDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
