package panes

import (
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type paneDetailModel struct {
	pane      Pane
	connector database.Connector
	logger    common.Logger
}

func NewPaneDetailModel(connector database.Connector, logger common.Logger, pane *Pane) tea.Model {
	if pane == nil {
		firstPane, err := gorm.G[Pane](connector.DB).First(connector.Ctx)
		if err != nil {
			logger.Errorlogger.Printf("Error loading first pane: %v", err)
		}
		pane = &firstPane
	}
	return paneDetailModel{*pane, connector, logger}
}

func (m paneDetailModel) Init() tea.Cmd { return nil }

func (m paneDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

// TODO: inconsistent with other detail models
func (m paneDetailModel) View() string {
	var items []string
	items = append(items, common.TitleStyle.Foreground(lipgloss.Color("2")).Render("Pane Details"))
	items = append(items, fmt.Sprintf("Name: %s", m.pane.Name))
	items = append(items, fmt.Sprintf("PWD: %s", m.pane.WorkingDirectory))
	var direction string
	switch m.pane.SplitDirection {
	case Horizontal:
		direction = "Horizontal"
	case Vertical:
		direction = "Vertical"
	default:
		direction = "Invalid direction"
	}
	items = append(items, fmt.Sprintf("Split direction: %s", direction))
	items = append(items, fmt.Sprintf("Split ratio: %d", m.pane.SplitRatio))
	items = append(items, fmt.Sprintf("Init cmd: %s", m.pane.ShellCmd))
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
