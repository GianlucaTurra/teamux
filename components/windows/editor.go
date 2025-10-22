package windows

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

const (
	creating = iota
	editing
	quitting
)

type WindowEditorModel struct {
	focusedIndex int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	mode         int
	window       *data.Window
	error        error
	db           *sql.DB
	logger       common.Logger
}

func NewWindowEditorModel(db *sql.DB, logger common.Logger) common.TeamuxModel {
	m := WindowEditorModel{
		inputs: make([]textinput.Model, 2),
		db:     db,
		logger: logger,
		mode:   creating,
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = common.FocusedStyle
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

func (m WindowEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WindowEditorModel) Update(msg tea.Msg) (common.TeamuxModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.InputErrMsg:
		m.error = msg.Err
		return m, nil
	case common.EditWMsg:
		m.mode = editing
		m.window = &msg.Window
		m.inputs[0].SetValue(msg.Window.Name)
		if msg.Window.WorkingDirectory == "$HOME" {
			m.inputs[1].SetValue("")
		} else {
			m.inputs[1].SetValue(msg.Window.WorkingDirectory)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewWindowEditorModel(m.db, m.logger), common.Browse
		case "ctrl+c":
			m.mode = quitting
			return m, common.Quit
		case "tab", "shift+tab", "up", "down":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusedIndex--
			} else {
				m.focusedIndex++
			}
			if m.focusedIndex >= len(m.inputs) {
				m.focusedIndex = 0
			}
			if m.focusedIndex <= 0 {
				m.focusedIndex = len(m.inputs)
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
				cmd = m.createWindow()
			case editing:
				cmd = m.editWindow()
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

func (m WindowEditorModel) View() string {
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
	// fmt.Fprintf(&b, "%s", m.help.View())
	return b.String()
}

func (m *WindowEditorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *WindowEditorModel) createWindow() tea.Cmd {
	window := data.NewWindow(m.inputs[0].Value(), m.inputs[1].Value(), m.db)
	if err := window.Save(); err != nil {
		m.error = err
		m.logger.Errorlogger.Printf("Error saving window: %v", err)
		return func() tea.Msg { return common.InputErrMsg{Err: err} }
	}
	m.error = nil
	return common.WindowCreated
}

func (m *WindowEditorModel) editWindow() tea.Cmd {
	m.window.Name = m.inputs[0].Value()
	m.window.WorkingDirectory = m.inputs[1].Value()
	if err := m.window.Save(); err != nil {
		m.error = err
		m.logger.Errorlogger.Printf("Error saving window: %v", err)
		return func() tea.Msg { return common.InputErrMsg{Err: err} }
	}
	m.error = nil
	// TODO: this is a little confusing
	return common.WindowCreated
}

func (m WindowEditorModel) GetDB() *sql.DB {
	return m.db
}

func (m WindowEditorModel) GetLogger() common.Logger {
	return m.logger
}
