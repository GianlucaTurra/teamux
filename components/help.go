package components

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type basicHelpKeymap struct {
	help key.Binding
	quit key.Binding
}

func (k basicHelpKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k basicHelpKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.help}, {k.quit},
	}
}

var basicHelpKeys = basicHelpKeymap{
	help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "show help")),
	quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
}

type HelpModel struct {
	basicKeys basicHelpKeymap
	help      help.Model
	model     common.HelpModel
}

func NewHelpModel(model common.HelpModel) HelpModel {
	return HelpModel{
		basicKeys: basicHelpKeys,
		help:      help.New(),
		model:     model,
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.ClearHelpMsg:
		switch msg.Tab {
		case common.SessionsContainer:
			m.model = sessions.NewSessionBrowserHelpModel()
		case common.WindwowsContainer:
			m.model = windows.NewWindowsBrowserHelpModel()
		case common.PanesContainer:
			m.model = panes.NewPanesBrowserHelpModel()
		}
		m.model.HideHelp()
	case common.ShowFullHelpMsg:
		m.model.ToggleHelp()
	}
	return m, nil
}

func (m HelpModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.help.View(m.basicKeys),
		m.model.View(),
	)
}

func (m *HelpModel) SetModel(model common.HelpModel) {
	m.model = model
}
