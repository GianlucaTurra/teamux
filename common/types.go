package common

type State int

const (
	Browsing State = iota
	Deleting
	Quitting
)

type Severity int

const (
	Info Severity = iota
	Warning
	Error
)

type ComponentWithHelp int

const (
	PaneBrowser ComponentWithHelp = iota
	PaneEditor
	WindowBrowser
	WindowEditor
	SessionBrowser
	SessionEditor
)

type FocusedTab int

const (
	SessionsContainer = iota
	WindwowsContainer
	PanesContainer
)
