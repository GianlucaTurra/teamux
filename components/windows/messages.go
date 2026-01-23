package windows

type (
	EditWMsg          struct{ Window Window }
	NewWFocus         struct{ Window Window }
	AssociatePanesMsg struct{ Window Window }
	EditWindowMsg     struct{ Window Window }
	WindowCreatedMsg  struct{}
)
