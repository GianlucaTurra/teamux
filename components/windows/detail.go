package windows

import (
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type WindowDetailModel struct {
	window    Window
	connector database.Connector
	logger    common.Logger
}

func NewWindowDetailModel(connector database.Connector, logger common.Logger, window *Window) WindowDetailModel {
	if window == nil {
		firstWindow, err := gorm.G[Window](connector.DB).First(connector.Ctx)
		if err != nil {
			logger.Errorlogger.Printf("Error loading first window, falling back to default one.\n %v", err)
		}
		window = &firstWindow
	}
	// TODO: should not be needed but needs tests
	// if err := window.GetAllPanes(); err != nil {
	// 	logger.Errorlogger.Printf("Error loading panes for first window.\n %v", err)
	// }
	return WindowDetailModel{*window, connector, logger}
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

func (m WindowDetailModel) Init() tea.Cmd { return nil }

func (m WindowDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

// TODO: should be a common function
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
