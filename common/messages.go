package common

type (
	ShowFullHelpMsg      struct{ Component ComponentWithHelp }
	ClearHelpMsg         struct{ Tab FocusedTab }
	InputErrMsg          struct{ Err error }
	OpenMsg              struct{}
	DeleteMsg            struct{}
	SwitchMsg            struct{}
	NewWindowMsg         struct{}
	NewPaneMsg           struct{}
	KillMsg              struct{}
	ReloadMsg            struct{}
	BrowseMsg            struct{}
	QuitMsg              struct{}
	NextTabMsg           struct{}
	PreviousTabMsg       struct{}
	UpDownMsg            struct{}
	LoadDataMsg          struct{}
	UpdateDetailMsg      struct{}
	CreateWindowMsg      struct{}
	SetOutputMsgTimerMsg struct{}
	ResetOutputMsgMsg    struct{}
	OutputMsg            struct {
		Err      error
		Severity Severity
	}
)
