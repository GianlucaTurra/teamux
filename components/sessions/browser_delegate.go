package sessions

import (
	"fmt"
	"io"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionDelegate struct{}

func (d sessionDelegate) Height() int                             { return 1 }
func (d sessionDelegate) Spacing() int                            { return 0 }
func (d sessionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d sessionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	fn := common.ItemStyle.Render
	if i.open {
		fn = func(s ...string) string { return common.OpenStyle.Render("* " + strings.Join(s, " ")) }
	}
	if index == m.Index() {
		if i.open {
			fn = func(s ...string) string {
				return common.SelectedOpenStyle.Render(">*" + strings.Join(s, " "))
			}
		} else {
			fn = func(s ...string) string {
				return common.SelectedStyle.Render("> " + strings.Join(s, " "))
			}
		}
	}
	fmt.Fprint(w, fn(str))
}
