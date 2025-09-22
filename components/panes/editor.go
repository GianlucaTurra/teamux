package panes

import (
	"database/sql"
	"fmt"
	"strconv"
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

type PaneEditorModel struct {
	focusedIndex int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	mode         int
	pane         *data.Pane
	db           *sql.DB
	logger       common.Logger
}

func NewPaneEditorModel(db *sql.DB, logger common.Logger) PaneEditorModel {
	m := PaneEditorModel{
		inputs: make([]textinput.Model, 4),
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
		case 2:
			t.Prompt = "Direction: "
			t.PromptStyle = common.BlurredStyle
			t.CharLimit = 1
		case 3:
			t.Prompt = "Ratio: "
			t.PromptStyle = common.BlurredStyle
			t.CharLimit = 3
		}
		m.inputs[i] = t
	}
	return m
}

func (m PaneEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PaneEditorModel) Update(msg tea.Msg) (PaneEditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.EditPMsg:
		m.mode = editing
		m.pane = &msg.Pane
		m.inputs[0].SetValue(msg.Pane.Name)
		if msg.Pane.WorkingDirectory == "$HOME" {
			m.inputs[1].SetValue("")
		} else {
			m.inputs[1].SetValue(msg.Pane.WorkingDirectory)
		}
		if msg.Pane.IsHorizontal() {
			m.inputs[2].SetValue("h")
		} else {
			m.inputs[2].SetValue("v")
		}
		m.inputs[3].SetValue(strconv.Itoa(msg.Pane.SplitRatio))
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewPaneEditorModel(m.db, m.logger), common.Browse
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
				cmd = m.createPane()
			case editing:
				cmd = m.editPane()
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

func (m PaneEditorModel) View() string {
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
	return b.String()
}

func (m *PaneEditorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *PaneEditorModel) createPane() tea.Cmd {
	var pane data.Pane
	ratio, err := strconv.Atoi(m.inputs[3].Value())
	if err != nil {
		m.logger.Errorlogger.Printf("Error converting ratio to int: %v", err)
	}
	switch strings.ToLower(m.inputs[2].Value()) {
	case "h":
		pane = data.NewHorizontalPane(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			ratio,
			m.db,
		)
	case "v":
		pane = data.NewVerticalPane(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			ratio,
			m.db,
		)
	default:
		err := fmt.Errorf("invalid direction: %s", m.inputs[2].Value())
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	if err := pane.Save(); err != nil {
		m.logger.Errorlogger.Printf("Error saving pane: %v", err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	return common.PaneCreated
}

func (m *PaneEditorModel) editPane() tea.Cmd {
	switch strings.ToLower(m.inputs[2].Value()) {
	case "h":
		m.pane.SetHorizontal()
	case "v":
		m.pane.SetVertical()
	default:
		err := fmt.Errorf("invalid direction: %s", m.inputs[2].Value())
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	ratio, err := strconv.Atoi(m.inputs[3].Value())
	if err != nil {
		m.logger.Errorlogger.Printf("Error converting ratio to int: %v", err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	m.pane.Name = m.inputs[0].Value()
	m.pane.WorkingDirectory = m.inputs[1].Value()
	m.pane.SplitRatio = ratio
	if err := m.pane.Save(); err != nil {
		m.logger.Errorlogger.Printf("Error saving pane: %v", err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	// TODO: kinda confusing
	return common.PaneCreated
}
