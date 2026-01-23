package sessions

import tea "github.com/charmbracelet/bubbletea"

func EditS(s Session) tea.Msg { return EditSMsg{s} }
func SessionCreated() tea.Msg { return SessionCreatedMsg{} }
func NewSession() tea.Msg     { return NewSessionMsg{} }
