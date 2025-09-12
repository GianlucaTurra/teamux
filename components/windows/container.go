package windows

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	tea "github.com/charmbracelet/bubbletea"
)

type WindowContainerModel struct {
	windowBrowser WindowBrowserModel
	windowEditor  WindowEditorModel
	focusedModel  int
}

const (
	windowBrowser = iota
	windowEditor
)

func NewWindowContainerModel(db *sql.DB, logger common.Logger) WindowContainerModel {
	return WindowContainerModel{
		windowBrowser: NewWindowBrowserModel(db, logger),
		windowEditor:  NewWindowEditorModel(db, logger),
		focusedModel:  0,
	}
}

func (m WindowContainerModel) Init() tea.Cmd {
	return tea.Batch(m.windowBrowser.Init())
}

func (m WindowContainerModel) Update(msg tea.Msg) (WindowContainerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewWindowMsg:
		m.focusedModel = windowEditor
		return m, nil
	case common.WindowCreatedMsg:
		m.focusedModel = windowBrowser
		return m, common.Reaload
	case common.BrowseMsg:
		m.focusedModel = windowBrowser
		return m, nil
	case common.EditWMsg:
		m.focusedModel = windowEditor
	case tea.KeyMsg:
		if msg.String() == "n" && m.focusedModel == windowBrowser && m.windowBrowser.state != common.Deleting {
			return m, common.NewWindow
		}
	}
	var cmds []tea.Cmd
	switch m.focusedModel {
	case windowEditor:
		newEditor, cmd := m.windowEditor.Update(msg)
		m.windowEditor = newEditor
		cmds = append(cmds, cmd)
	case windowBrowser:
		newBrowser, cmd := m.windowBrowser.Update(msg)
		m.windowBrowser = newBrowser
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m WindowContainerModel) View() string {
	switch m.focusedModel {
	case windowBrowser:
		return m.windowBrowser.View()
	case windowEditor:
		return m.windowEditor.View()
	default:
		return ""
	}
}
