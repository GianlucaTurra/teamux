package sessions

import (
	"database/sql"
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SessionDetailModel struct {
	session data.Session
	db      *sql.DB
	logger  common.Logger
}

func NewSessionTreeModel(db *sql.DB, logger common.Logger, session *data.Session) SessionDetailModel {
	if session == nil {
		firstSession, err := data.GetFirstSession(db)
		if err != nil {
			logger.Errorlogger.Printf("Error loading first session, falling back to default one.\n %v", err)
		}
		if err := firstSession.GetPWD(); err != nil {
			logger.Errorlogger.Printf("Error loading working directory for first session, falling back to blank directory.\n %v", err)
		}
		session = &firstSession
	}
	if err := session.GetAllWindows(); err != nil {
		logger.Errorlogger.Printf("Error loading windows for session %s\n%v", session.Name, err)
	}
	return SessionDetailModel{*session, db, logger}
}

func (m SessionDetailModel) View() string {
	var items []string
	title := renderTreeItem("Session Details", "", 0, false)
	items = append(items, common.TitleStyle.Foreground(lipgloss.Color("2")).Render(title))
	items = append(items, renderTreeItem(m.session.Name, m.session.WorkingDirectory, 0, false))
	for i, w := range m.session.Windows {
		items = append(items, renderTreeItem(w.Name, w.WorkingDirectory, 1, i == len(m.session.Windows)-1))
		if err := w.GetAllPanes(); err != nil {
			m.logger.Errorlogger.Printf("Error loading panes for window %s\n%v", w.Name, err)
			continue
		}
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
		// prefix = "│   " + treeSymbol + " "
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
