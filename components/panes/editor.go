package panes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/database"
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
	mode         int
	pane         *Pane
	connector    database.Connector
	logger       common.Logger
}

func NewPaneEditorModel(connector database.Connector, logger common.Logger, pane *Pane) PaneEditorModel {
	m := PaneEditorModel{
		inputs:    make([]textinput.Model, 6),
		connector: connector,
		logger:    logger,
		mode:      creating,
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
			t.Prompt = "Direction (v/h): "
			t.PromptStyle = common.BlurredStyle
			t.CharLimit = 1
		case 3:
			t.Prompt = "Ratio: "
			t.PromptStyle = common.BlurredStyle
			t.CharLimit = 3
		case 4:
			t.Prompt = "Target: "
			t.PromptStyle = common.BlurredStyle
			t.CharLimit = 50
		case 5:
			t.Prompt = "ShellCmd: "
			t.PromptStyle = common.BlurredStyle
			t.CharLimit = 100
		}
		m.inputs[i] = t
	}
	if pane != nil {
		m.mode = editing
		m.inputs[0].SetValue(pane.Name)
		m.inputs[1].SetValue(pane.WorkingDirectory)
		m.inputs[2].SetValue(pane.SplitDirection)
		m.inputs[3].SetValue(strconv.Itoa(pane.SplitRatio))
		m.inputs[4].SetValue(pane.Target)
		m.inputs[5].SetValue(pane.ShellCmd)
	}
	return m
}

func (m PaneEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PaneEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case EditPMsg:
		m.pane = &msg.Pane
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewPaneEditorModel(m.connector, m.logger, nil), common.Browse
		case "?":
			return m, func() tea.Msg { return common.ShowFullHelpMsg{Component: common.PaneEditor} }
		case "ctrl+c":
			m.mode = quitting
			return m, common.Quit
		case "tab", "shift+tab", "up", "down":
			return m.cycleInputs(msg.String())
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

// FIXME: should be a tea.Cmd itself
func (m *PaneEditorModel) createPane() tea.Cmd {
	ratio, err := strconv.Atoi(m.inputs[3].Value())
	if err != nil {
		m.logger.Errorlogger.Printf("Error converting ratio to int: %v", err)
	}
	switch strings.ToLower(m.inputs[2].Value()) {
	case "h":
		_, err = CreateHorizontalPane(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			ratio,
			m.connector,
			m.inputs[4].Value(),
			m.inputs[5].Value(),
		)
	case "v":
		_, err = CreateVerticalPane(
			m.inputs[0].Value(),
			m.inputs[1].Value(),
			ratio,
			m.connector,
			m.inputs[4].Value(),
			m.inputs[5].Value(),
		)
	default:
		err := fmt.Errorf("invalid direction: %s", m.inputs[2].Value())
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	if err != nil {
		m.logger.Errorlogger.Printf("Error saving pane: %v", err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	return func() tea.Msg { return PanesEditedMsg{} }
}

// FIXME: should be a tea.Cmd itself
func (m *PaneEditorModel) editPane() tea.Cmd {
	m.pane.Name = m.inputs[0].Value()
	m.pane.WorkingDirectory = m.inputs[1].Value()
	m.pane.SplitDirection = m.inputs[2].Value()
	ratio, err := strconv.Atoi(m.inputs[3].Value())
	if err != nil {
		m.logger.Errorlogger.Printf("Error converting ratio to int: %v", err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	m.pane.SplitRatio = ratio
	m.pane.Target = m.inputs[4].Value()
	if _, err := m.pane.Save(m.connector); err != nil {
		m.logger.Errorlogger.Printf("Error saving pane: %v", err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	return func() tea.Msg { return PanesEditedMsg{} }
}

func (m *PaneEditorModel) cycleInputs(s string) (tea.Model, tea.Cmd) {
	if s == "up" || s == "shift+tab" {
		m.focusedIndex--
	} else {
		m.focusedIndex++
	}
	if m.focusedIndex == len(m.inputs) {
		m.focusedIndex = 0
	}
	if m.focusedIndex == -1 {
		m.focusedIndex = len(m.inputs) - 1
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
}
