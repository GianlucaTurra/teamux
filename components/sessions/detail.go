package sessions

import (
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type SessionDetailModel struct {
	session   data.Session
	connector data.Connector
	logger    common.Logger
}

func NewSessionTreeModel(connector data.Connector, logger common.Logger, session *data.Session) SessionDetailModel {
	if session == nil {
		firstSession, err := gorm.G[data.Session](connector.DB).First(connector.Ctx)
		if err != nil {
			logger.Errorlogger.Printf("Error loading first session, falling back to default one.\n %v", err)
		}
		session = &firstSession
	}
	return SessionDetailModel{*session, connector, logger}
}

func (m SessionDetailModel) View() string {
	var items []string
	title := renderTreeItem("Session Details", "", 0, false)
	items = append(items, common.TitleStyle.Foreground(lipgloss.Color("2")).Render(title))
	items = append(items, renderTreeItem(m.session.Name, m.session.WorkingDirectory, 0, false))
	for i, w := range m.session.Windows {
		items = append(items, renderTreeItem(w.Name, w.WorkingDirectory, 1, i == len(m.session.Windows)-1))
		for j, p := range w.Panes {
			items = append(items, renderTreeItem(p.Name, p.WorkingDirectory, 2, j == len(w.Panes)-1))
		}
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
	case 2:
		prefix = "   " + treeSymbol + " "
	}
	nameCol := fmt.Sprintf("%-*s", 16, prefix+name)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		common.ItemStyle.Render(nameCol),
		common.ItemStyle.Italic(true).Render(pwd),
	)
}

func (m SessionDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SessionDetailModel) Init() tea.Cmd {
	return nil
}
