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
	sessionTreeModel struct {
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
	padding := 1 * (i.level * 2)
	fn := common.ItemStyle.PaddingLeft(padding).Render
	str := fmt.Sprintf("%s %s", treeSymbol, i.title)
	fmt.Fprint(w, fn(str))
}

func (s treeItem) FilterValue() string { return "" }

func NewSessionTreeModel(db *sql.DB, logger common.Logger, session data.Session) {}
