package common

import "github.com/GianlucaTurra/teamux/components/data"

type (
	TmuxSessionsChanged struct{}
	OpenMsg             struct{}
	DeleteMsg           struct{}
	SwitchMsg           struct{}
	NewSessionMsg       struct{}
	NewWindowMsg        struct{}
	NewPaneMsg          struct{}
	KillMsg             struct{}
	InputErrMsg         struct{ Err error }
	ReloadMsg           struct{}
	SessionCreatedMsg   struct{}
	WindowCreatedMsg    struct{}
	PanesEditedMsg      struct{}
	BrowseMsg           struct{}
	EditSMsg            struct{ Session data.Session }
	EditWMsg            struct{ Window data.Window }
	EditPMsg            struct{ Pane data.Pane }
	QuitMsg             struct{}
	NextTabMsg          struct{}
	PreviousTabMsg      struct{}
	UpDownMsg           struct{}
	NewSFocus           struct{ Session data.Session }
	NewWFocus           struct{ Window data.Window }
	OutputMsg           struct {
		Err      error
		Severity Severity
	}
	ShowFullHelpMsg struct{ Component ComponentWithHelp }
	ClearHelpMsg    struct{ Tab FocusedTab }
)
