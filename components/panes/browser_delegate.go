package panes

import (
	"fmt"
	"io"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type paneDelegate struct{}

func (d paneDelegate) Height() int                             { return 1 }
func (d paneDelegate) Spacing() int                            { return 0 }
func (d paneDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d paneDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(paneItem)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	fn := common.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return common.SelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}
