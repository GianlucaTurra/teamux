package windows

import (
	"fmt"
	"io"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: since it's common for all list just move it to the common package
type windowPanesDelegate struct{}

func (d windowPanesDelegate) Height() int                             { return 1 }
func (d windowPanesDelegate) Spacing() int                            { return 0 }
func (d windowPanesDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d windowPanesDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(availablePanes)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	fn := common.ItemStyle.Render
	if i.selected {
		fn = func(s ...string) string { return common.OpenStyle.Render(" ✓" + strings.Join(s, " ")) }
	}
	if index == m.Index() {
		if i.selected {
			fn = func(s ...string) string {
				return common.SelectedOpenStyle.Render(">✓" + strings.Join(s, " "))
			}
		} else {
			fn = func(s ...string) string {
				return common.SelectedStyle.Render("> " + strings.Join(s, " "))
			}
		}
	}
	fmt.Fprint(w, fn(str))
}
