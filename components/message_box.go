package components

import (
	"github.com/GianlucaTurra/teamux/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MessageBoxModel struct {
	message  string
	severity common.Severity
}

func NewMessageBoxModel() tea.Model {
	return MessageBoxModel{"", common.Info}
}

func (m MessageBoxModel) Init() tea.Cmd {
	return nil
}

func (m MessageBoxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.ResetOutputMsgMsg:
		m.message = ""
		m.severity = common.Info
		return m, nil
	case common.OutputMsg:
		m.message = msg.Err.Error()
		m.severity = msg.Severity
		return m, func() tea.Msg { return common.SetOutputMsgTimerMsg{} }
	default:
		return m, nil
	}
}

func (m MessageBoxModel) View() string {
	var color string
	switch m.severity {
	case common.Info:
		color = common.White
	case common.Warning:
		color = common.Yellow
	case common.Error:
		color = common.Red
	default:
		color = common.White
	}
	return common.ItemStyle.Foreground(lipgloss.Color(color)).Render(m.message)
}
