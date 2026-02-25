package data

import (
	"testing"

	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/tmux"
)

func TestWindow_openAndCascade(t *testing.T) {
	testSession := "Test"
	if err := tmux.NewSession(testSession, "~/"); err != nil {
		t.Errorf("unable to open test sessions: %v", err)
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		target  *string
		w       windows.Window
		wantErr bool
	}{
		{
			name:    "Complete target",
			target:  &testSession,
			w:       windows.Window{Name: "TestW", WorkingDirectory: "~/"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// FIXME: non public field
			// p := panes.Pane{Name: "TestP", WorkingDirectory: "~/", SplitRatio: 25, splitDirection: 1}
			// p2 := panes.Pane{Name: "TestP2", WorkingDirectory: "~/", SplitRatio: 50, splitDirection: 0}
			// tt.w.Panes = append(tt.w.Panes, p, p2)
			// gotErr := tt.w.openAndCascade(tt.target)
			// if gotErr != nil {
			// 	if !tt.wantErr {
			// 		t.Errorf("openAndCascade() failed: %v", gotErr)
			// 	}
			// 	return
			// }
			// if tt.wantErr {
			// 	t.Fatal("openAndCascade() succeeded unexpectedly")
			// }
		})
	}
	if err := tmux.KillSession(testSession); err != nil {
		t.Errorf("unable to kill test session: %v", err)
	}
}
