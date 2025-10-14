package sessions

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var cursorStyle = common.FocusedStyle

const (
	creating = iota
	editing
	quitting
)

type SessionEditorModel struct {
	focusedIndex int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	mode         int
	session      *data.Session
	error        error
	db           *sql.DB
	logger       common.Logger
}

func NewSessionEditorModel(db *sql.DB, logger common.Logger) SessionEditorModel {
	m := SessionEditorModel{
		inputs: make([]textinput.Model, 2),
		db:     db,
		logger: logger,
		mode:   creating,
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Prompt = "Name: "
			t.PromptStyle = common.BlurredStyle
			t.Focus()
			t.PromptStyle = common.FocusedStyle
			t.TextStyle = common.FocusedStyle
			t.CharLimit = 50
		case 1:
			t.Prompt = "WorkDir: "
			t.PromptStyle = common.BlurredStyle
			t.CharLimit = 100
		}
		m.inputs[i] = t
	}
	return m
}

func (m SessionEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SessionEditorModel) Update(msg tea.Msg) (common.TeamuxModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.InputErrMsg:
		m.error = msg.Err
		return m, nil
	case common.EditSMsg:
		m.mode = editing
		m.session = &msg.Session
		m.inputs[0].SetValue(msg.Session.Name)
		if msg.Session.WorkingDirectory == "$HOME" {
			m.inputs[1].SetValue("")
		} else {
			m.inputs[1].SetValue(msg.Session.WorkingDirectory)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewSessionEditorModel(m.db, m.logger), common.Browse
		case "ctrl+c":
			m.mode = quitting
			return m, common.Quit
		case "tab", "shift+tab":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusedIndex--
			} else {
				m.focusedIndex++
			}
			if m.focusedIndex == -1 {
				m.focusedIndex = len(m.inputs) - 1
			}
			if m.focusedIndex == len(m.inputs) {
				m.focusedIndex = 0
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusedIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = common.FocusedStyle
					m.inputs[i].TextStyle = common.FocusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = common.BlurredStyle
				m.inputs[i].TextStyle = common.BlurredStyle
			}
			return m, tea.Batch(cmds...)
		case "enter":
			var cmd tea.Cmd
			switch m.mode {
			case creating:
				cmd = m.createSession()
			case editing:
				cmd = m.editSession()
			}
			m.focusedIndex = 0
			for i := range m.inputs {
				m.inputs[i].Reset()
			}
			return m, cmd
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
		m.error = err
		m.logger.Errorlogger.Printf("Error saving session: %v", err)
		return func() tea.Msg { return common.InputErrMsg{Err: err} }
	}
	m.error = nil
	return common.SessionCreated
}

func (m *SessionEditorModel) editSession() tea.Cmd {
	m.session.Name = m.inputs[0].Value()
	m.session.WorkingDirectory = m.inputs[1].Value()
	if err := m.session.Save(); err != nil {
		m.error = err
		m.logger.Errorlogger.Printf("Error saving session: %v", err)
		return func() tea.Msg { return common.InputErrMsg{Err: err} }
	}
	m.error = nil
	// TODO: this is a little confusing
	return common.SessionCreated
}

func (m SessionEditorModel) View() string {
	if m.mode == quitting {
		return ""
	}
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	fmt.Fprintf(&b, "\n\n")
	if m.error != nil {
		fmt.Fprintf(&b, "\nError: %v", m.error)
	}
	return b.String()
}

func (m SessionEditorModel) GetDB() *sql.DB {
	return m.db
}

func (m SessionEditorModel) GetLogger() common.Logger {
	return m.logger
}
