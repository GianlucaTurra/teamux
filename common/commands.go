package common

import (
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
)

func Open() tea.Msg                  { return OpenMsg{} }
func Delete() tea.Msg                { return DeleteMsg{} }
func Switch() tea.Msg                { return SwitchMsg{} }
func NewSession() tea.Msg            { return NewSessionMsg{} }
func NewWindow() tea.Msg             { return NewWindowMsg{} }
func NewPane() tea.Msg               { return NewPaneMsg{} }
func Kill() tea.Msg                  { return KillMsg{} }
func SessionCreated() tea.Msg        { return SessionCreatedMsg{} }
func WindowCreated() tea.Msg         { return WindowCreatedMsg{} }
func PaneEdited() tea.Msg            { return PanesEditedMsg{} }
func Reaload() tea.Msg               { return ReloadMsg{} }
func Browse() tea.Msg                { return BrowseMsg{} }
func EditS(s data.Session) tea.Msg   { return EditSMsg{s} }
func EditW(w data.Window) tea.Msg    { return EditWMsg{w} }
func EditP(p data.Pane) tea.Msg      { return EditPMsg{p} }
func Quit() tea.Msg                  { return QuitMsg{} }
func NextTab() tea.Msg               { return NextTabMsg{} }
func PreviousTab() tea.Msg           { return PreviousTabMsg{} }
func UpDown() tea.Msg                { return UpDownMsg{} }
func ClearHelp(t FocusedTab) tea.Msg { return ClearHelpMsg{t} }
func LoadData() tea.Msg              { return LoadDataMsg{} }
func UpdateDetail() tea.Msg          { return UpdateDetailMsg{} }
