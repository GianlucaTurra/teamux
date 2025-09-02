package sessions

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	treeDelegate struct{}
	treeItem     struct {
		title            string
		workingDirectory string
		level            int
		isLast           bool
	}
	SessionTreeModel struct {
		list   list.Model
		db     *sql.DB
		logger common.Logger
	}
)

func (d treeDelegate) Height() int                             { return 1 }
func (d treeDelegate) Spacing() int                            { return 0 }
func (d treeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d treeDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(treeItem)
	if !ok {
		return
	}
	var treeSymbol string
	if i.level != 0 {
		if i.isLast {
			treeSymbol = "└──"
		} else {
			treeSymbol = "├──"
		}
	}
	padding := 1 * (i.level * 2)
	fn := common.ItemStyle.PaddingLeft(padding).Render
	fnitalics := common.ItemStyle.Italic(true).Render
	str := fmt.Sprintf("%s %s", treeSymbol, i.title)
	fmt.Fprint(w, fn(str), fnitalics(i.workingDirectory))
}

func (s treeItem) FilterValue() string { return "" }

func NewSessionTreeModel(db *sql.DB, logger common.Logger, session *data.Session) SessionTreeModel {
	if session == nil {
		firstSession, err := data.GetFirstSession(db)
		if err != nil {
			logger.Errorlogger.Printf("Error loading first session, falling back to default one.\n %v", err)
		}
		session = &firstSession
	}
	if err := session.GetAllWindows(); err != nil {
		logger.Errorlogger.Printf("Error loading windows for session %s\n%v", session.Name, err)
	}
	layouts := []list.Item{}
	layouts = append(layouts, treeItem{session.Name, session.WorkingDirectory, 0, false})
	for i, window := range session.Windows {
		isLast := false
		if i == len(session.Windows) {
			isLast = true
		}
		layouts = append(layouts, treeItem{window.Name, window.WorkingDirectory, 1, isLast})
	}
	l := list.New(layouts, treeDelegate{}, 100, 10)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	return SessionTreeModel{l, db, logger}
}

func (m SessionTreeModel) View() string {
	return m.list.View()
}

func (m SessionTreeModel) Update(msg tea.Msg) (SessionTreeModel, tea.Cmd) {
	return m, nil
}

func (m SessionTreeModel) Init() tea.Cmd {
	return nil
}
