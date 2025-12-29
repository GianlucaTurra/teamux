package components

import (
	"github.com/GianlucaTurra/teamux/common"
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

func NewHelpModel() HelpModel {
	return HelpModel{
		basicKeys: basicHelpKeys,
		help:      help.New(),
		model:     NewBrowserHelpModel(),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg.(type) {
	case common.ClearHelpMsg:
		m.model = NewBrowserHelpModel()
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
