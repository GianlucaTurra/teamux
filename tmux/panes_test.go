package tmux_test

import (
	"testing"

	"github.com/GianlucaTurra/teamux/tmux"
)

func TestSplitWindowWithTargetWindow(t *testing.T) {
	tests := []struct {
		name string
		// Named input parameters for target function.
		target           string
		splitRatio       int
		workingDirectory string
		horizontal       bool
		wantErr          bool
		currentSession   bool
	}{
		{
			name:             "Empty directory pane with target",
			target:           "test",
			splitRatio:       25,
			workingDirectory: "",
			horizontal:       true,
			wantErr:          false,
			currentSession:   true,
		},
		{
			name:             "Empty directory pane without target",
			target:           "",
			splitRatio:       25,
			workingDirectory: "",
			horizontal:       true,
			wantErr:          true,
			currentSession:   true,
		},
		{
			name:             "Directory with target",
			target:           "test",
			splitRatio:       25,
			workingDirectory: "~/go/",
			horizontal:       true,
			wantErr:          false,
			currentSession:   true,
		},
		{
			name:             "Directory with target and no splitRatio",
			target:           "test",
			workingDirectory: "~/go/",
			horizontal:       true,
			wantErr:          true,
			currentSession:   true,
		},
		{
			name:             "Detached session",
			target:           "Test:test",
			workingDirectory: "",
			horizontal:       true,
			splitRatio:       25,
			wantErr:          false,
			currentSession:   false,
		},
		{
			name:             "Non existing target",
			target:           "i do not exist",
			workingDirectory: "",
			horizontal:       true,
			splitRatio:       25,
			wantErr:          true,
			currentSession:   true,
		},
	}
	for _, tt := range tests {
		if tt.currentSession {
			if err := tmux.NewWindow("test", "", nil); err != nil {
				t.Errorf("Error opening test window: %v", err)
			}
		} else {
			if err := tmux.NewSession("Test", ""); err != nil {
				t.Errorf("Error opening test session: %v", err)
			}
			target := "Test"
			if err := tmux.NewWindow("test", "", &target); err != nil {
				t.Errorf("Error opening test window: %v", err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tmux.SplitWindowWithTargetWindow(tt.target, tt.splitRatio, tt.workingDirectory, tt.horizontal)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SplitWindowWithTarget() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SplitWindowWithTarget() succeeded unexpectedly")
			}
		})
		if tt.currentSession {
			if err := tmux.KillWindow("test"); err != nil {
				t.Errorf("Error killing test window: %v", err)
			}
		} else {
			if err := tmux.KillSession("Test"); err != nil {
				t.Errorf("Error killing test session: %v", err)
			}
		}
	}
}
