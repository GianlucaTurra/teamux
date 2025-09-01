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
		title  string
		level  int
		isLast bool
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
	if i.isLast {
		treeSymbol = "└──"
	} else {
		treeSymbol = "├──"
	}
	padding := i.level * 2
	fn := common.ItemStyle.PaddingLeft(padding).Render
	str := fmt.Sprintf("%s %s", treeSymbol, i.title)
	fmt.Fprint(w, fn(str))
}

func (s treeItem) FilterValue() string { return "" }

func NewSessionTreeModel(db *sql.DB, logger common.Logger, session data.Session) SessionTreeModel {
	layouts := []list.Item{}
	if err := session.GetAllWindows(); err != nil {
		logger.Errorlogger.Printf("Error loading sessions windows %v", err)
	}
	layouts = append(layouts, treeItem{session.Name, 0, false})
	for i, window := range session.Windows {
		if i == len(session.Windows) {
			layouts = append(layouts, treeItem{window.Name, 1, true})
		} else {
			layouts = append(layouts, treeItem{window.Name, 1, false})
		}
	}
	l := list.New(layouts, treeDelegate{}, 100, 10)
	return SessionTreeModel{l, db, logger}
}
