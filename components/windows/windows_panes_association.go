package windows

import (
	"errors"
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type availablePanes struct {
	pane     data.Pane
	title    string
	desc     string
	selected bool
}

func (ap availablePanes) FilterValue() string {
	return ""
}

type WindowPanesAssociationModel struct {
	model     list.Model
	connector data.Connector
	logger    common.Logger
	window    data.Window
	state     common.State
}

func NewWindowPanesAssociationModel(
	connector data.Connector,
	logger common.Logger,
	window data.Window,
) WindowPanesAssociationModel {
	var aps []list.Item
	l := list.New(aps, windowPanesDelegate{}, 100, 10)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(true)
	l.Styles.Title = common.HeaderStyle
	l.Title = "Associate panes to: " + window.Name
	l.Styles.PaginationStyle = common.PaginationStyle
	return WindowPanesAssociationModel{
		model:     l,
		connector: connector,
		logger:    logger,
		window:    window,
	}
}

func (m WindowPanesAssociationModel) Init() tea.Cmd {
	return nil
}

func (m WindowPanesAssociationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.LoadDataMsg:
		return m.loadData()
	case common.UpdateDetailMsg:
		return m, func() tea.Msg { return common.NewWFocus{Window: m.window} }
	case tea.KeyMsg:
		switch msg.String() {
		// TODO: handle keys to create and edit panes from here?
		case " ":
			return m.selectPane()
		case "esc":
			return m, common.Browse
		case "q", "ctrl+c":
			m.state = common.Quitting
			return m, common.Quit
		}
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m WindowPanesAssociationModel) View() string {
	switch m.state {
	case common.Quitting:
		return ""
	default:
		return m.model.View()
	}
}

func (m *WindowPanesAssociationModel) loadData() (tea.Model, tea.Cmd) {
	var aps []list.Item
	panes, err := data.ReadAllPanes(m.connector.DB)
	if err != nil {
		errMsg := fmt.Sprintf("error reading panes: %v", err)
		m.logger.Errorlogger.Println(errMsg)
		return m, func() tea.Msg { return common.OutputMsg{Err: errors.New(errMsg), Severity: common.Error} }
	}
	for _, pane := range panes {
		selected := false
		for _, p := range m.window.Panes {
			if pane.ID == p.ID {
				selected = true
				break
			}
		}
		aps = append(aps, availablePanes{pane: pane, title: pane.Name, selected: selected})
	}
	m.model.SetItems(aps)
	return m, nil
}

func (m *WindowPanesAssociationModel) selectPane() (tea.Model, tea.Cmd) {
	p := m.model.SelectedItem().(availablePanes)
	var err error
	if p.selected {
		err = m.connector.DB.Model(&m.window).Association("Panes").Delete(&p.pane)
	} else {
		err = m.connector.DB.Model(&m.window).Association("Panes").Append(&p.pane)
	}
	if err != nil {
		errMsg := fmt.Sprintf("error appending/deleting pane %s to window %s: %v", p.pane.Name, m.window.Name, err)
		m.logger.Errorlogger.Println(errMsg)
		// TODO: refactor OutputMsg with methods for clearer construction
		return m, func() tea.Msg { return common.OutputMsg{Err: errors.New(errMsg), Severity: common.Error} }
	}
	itemIndex := m.model.GlobalIndex()
	p.selected = !p.selected
	m.model.SetItem(itemIndex, p)
	return m, common.UpdateDetail
}
