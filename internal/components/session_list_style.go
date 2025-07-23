package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle           = lipgloss.NewStyle().MarginLeft(2)
	sessionStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedSessionStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle      = list.DefaultStyles().TitleBar.PaddingLeft(4)
	helpStyle            = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type SessionDelegate struct{}

func (d SessionDelegate) Height() int                             { return 1 }
func (d SessionDelegate) Spacing() int                            { return 0 }
func (d SessionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SessionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Session)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i)
	fn := sessionStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedSessionStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}
