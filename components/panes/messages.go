package panes

type (
	EditPMsg       struct{ Pane Pane }
	PaneCreatedMsg struct{}
	PanesEditedMsg struct{}
)
