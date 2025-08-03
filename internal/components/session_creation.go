package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle   = focusedStyle
	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type sessionInputModel struct {
	focusedIndex int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	quitting     bool
}
type sessionInputInfo struct {
	File string
	Name string
}

func newSessionInputModel() sessionInputModel {
	m := sessionInputModel{inputs: make([]textinput.Model, 2)}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Prompt = "File: "
			t.PromptStyle = blurredStyle
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 50
		case 1:
			t.Prompt = "Name: "
			t.PromptStyle = blurredStyle
			t.CharLimit = 20
		}
		m.inputs[i] = t
	}
	return m
}

func (m sessionInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m sessionInputModel) Update(msg tea.Msg) (sessionInputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab", "shift+tab", "up", "down":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusedIndex--
			} else {
				m.focusedIndex++
			}
			if m.focusedIndex > len(m.inputs) {
				m.focusedIndex = 0
			}
			if m.focusedIndex < 0 {
				m.focusedIndex = len(m.inputs)
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusedIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = blurredStyle
				m.inputs[i].TextStyle = blurredStyle
			}
			return m, tea.Batch(cmds...)
		case "enter":
			if m.focusedIndex == len(m.inputs) {
				m.readInputs()
				return m, tea.Quit
			}
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *sessionInputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *sessionInputModel) readInputs() {
	info := sessionInputInfo{
		File: m.inputs[0].Value(),
		Name: m.inputs[1].Value(),
	}
	fmt.Printf("Session Info: %+v\n", info)
}

func (m sessionInputModel) View() string {
	if m.quitting {
		return ""
	}
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.focusedIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	return b.String()
}
