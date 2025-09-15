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
