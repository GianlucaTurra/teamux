package common

import "github.com/GianlucaTurra/teamux/components/data"

type (
	TmuxSessionsChanged struct{}
	TmuxErr             struct{}
	OpenMsg             struct{}
	DeleteMsg           struct{}
	SwitchMsg           struct{}
	NewSessionMsg       struct{}
	NewWindowMsg        struct{}
	KillMsg             struct{}
	InputErrMsg         struct{ Err error }
	ReloadMsg           struct{}
	SessionCreatedMsg   struct{}
	BrowseMsg           struct{}
	EditMsg             struct{ Session data.Session }
	QuitMsg             struct{}
	NextTabMsg          struct{}
	PreviousTabMsg      struct{}
)
