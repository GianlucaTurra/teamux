package panes

type (
	EditPMsg       struct{ Pane Pane }
	NewPFocusMsg   struct{ Pane Pane }
	PaneCreatedMsg struct{}
	PanesEditedMsg struct{}
)
