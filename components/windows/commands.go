package windows

import tea "github.com/charmbracelet/bubbletea"

func EditW(w Window) tea.Msg { return EditWMsg{w} }
func WindowCreated() tea.Msg { return WindowCreatedMsg{} }
