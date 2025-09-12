package panes

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	tea "github.com/charmbracelet/bubbletea"
)

type PaneContainerModel struct {
	paneBrowser  PaneBrowserModel
	paneEditor   PaneEditorModel
	focusedModel int
}

const (
	paneBrowser = iota
	paneEditor
)

func NewPaneContainerModel(db *sql.DB, logger common.Logger) PaneContainerModel {
	return PaneContainerModel{
		paneBrowser:  NewPaneBrowserModel(db, logger),
		paneEditor:   NewPaneEditorModel(db, logger),
		focusedModel: 0,
	}
}

func (m PaneContainerModel) Init() tea.Cmd {
	return tea.Batch(m.paneBrowser.Init())
}

func (m PaneContainerModel) Update(msg tea.Msg) (PaneContainerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewPaneMsg:
		m.focusedModel = paneEditor
		return m, nil
	case common.PaneCreatedMsg:
		m.focusedModel = paneBrowser
		return m, common.Reaload
	case common.BrowseMsg:
		m.focusedModel = paneBrowser
		return m, nil
	case common.EditPMsg:
		m.focusedModel = paneEditor
	case tea.KeyMsg:
		if msg.String() == "n" && m.focusedModel == paneBrowser && m.paneBrowser.state != common.Deleting {
			return m, common.NewPane
		}
	}
	var cmds []tea.Cmd
	switch m.focusedModel {
	case paneEditor:
		newEditor, cmd := m.paneEditor.Update(msg)
		m.paneEditor = newEditor
		cmds = append(cmds, cmd)
	case paneBrowser:
		newBrowser, cmd := m.paneBrowser.Update(msg)
		m.paneBrowser = newBrowser
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m PaneContainerModel) View() string {
	switch m.focusedModel {
	case paneBrowser:
		return m.paneBrowser.View()
	case paneEditor:
		return m.paneEditor.View()
	default:
		return ""
	}
}
