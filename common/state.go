package common

type State int

const (
	Browsing State = iota
	Deleting
	Quitting
)
