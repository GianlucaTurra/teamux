package sessions

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
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

type SessionEditorModel struct {
	focusedIndex int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	quitting     bool
	error        error
	db           *sql.DB
	logger       common.Logger
	help         sessionEditorHelpModel
}

func NewSessionEditorModel(db *sql.DB, logger common.Logger) SessionEditorModel {
	m := SessionEditorModel{
		inputs: make([]textinput.Model, 2),
		db:     db,
		logger: logger,
		help:   newSessionEditorHelpModel(),
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Prompt = "Name: "
			t.PromptStyle = blurredStyle
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 50
		case 1:
			t.Prompt = "WorkDir: "
			t.PromptStyle = blurredStyle
			t.CharLimit = 100
		}
		m.inputs[i] = t
	}
	return m
}

func (m SessionEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SessionEditorModel) Update(msg tea.Msg) (SessionEditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.InputErrMsg:
		m.error = msg.Err
		return m, nil
	case common.EditMsg:
		m.inputs[0].SetValue(msg.Session.Name)
		var homeDir string
		var err error
		if homeDir, err = os.UserHomeDir(); err != nil {
			m.logger.Errorlogger.Printf("Error getting home directory: %v", err)
		}
		if homeDir == msg.Session.WorkingDirectory {
			m.inputs[1].SetValue("")
		} else {
			m.inputs[1].SetValue(msg.Session.WorkingDirectory)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			var cmd tea.Cmd
			m.help, cmd = m.help.Update(msg)
			return m, cmd
		case "esc":
			return m, common.Browse
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "tab", "shift+tab":
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
				m.focusedIndex = 0
				for i := range m.inputs {
					m.inputs[i].Reset()
				}
				return m, m.createSession()
			}
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *SessionEditorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *SessionEditorModel) createSession() tea.Cmd {
	session := data.NewSession(m.inputs[0].Value(), m.inputs[1].Value(), m.db)
	if err := session.Save(); err != nil {
		m.logger.Errorlogger.Printf("Error saving session: %v", err)
		return func() tea.Msg { return common.InputErrMsg{Err: err} }
	}
	m.error = nil
	return common.Created
}

func (m SessionEditorModel) View() string {
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
	if m.error != nil {
		fmt.Fprintf(&b, "\nError: %v", m.error)
	}
	fmt.Fprintf(&b, "%s", m.help.View())
	return b.String()
}
