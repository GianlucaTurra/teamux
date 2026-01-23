package panes

import tea "github.com/charmbracelet/bubbletea"

func PaneEdited() tea.Msg  { return PanesEditedMsg{} }
func EditP(p Pane) tea.Msg { return EditPMsg{p} }
