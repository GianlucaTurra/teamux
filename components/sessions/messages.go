package sessions

type (
	EditSMsg            struct{ Session Session }
	NewSFocus           struct{ Session Session }
	AssociateWindowsMsg struct{ Session Session }
	SessionCreatedMsg   struct{}
	NewSessionMsg       struct{}
	TmuxSessionsChanged struct{}
)
