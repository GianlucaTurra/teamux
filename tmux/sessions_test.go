package tmux_test

import (
	"testing"

	"github.com/GianlucaTurra/teamux/tmux"
)

func TestNewSession(t *testing.T) {
	tests := []struct {
		testName string
		// Named input parameters for target function.
		name             string
		workingDirectory string
		wantErr          bool
	}{
		{
			testName:         "Empty workingDirectory",
			name:             "Test",
			workingDirectory: "",
			wantErr:          false,
		},
		{
			testName:         "Home workingDirectory",
			name:             "Home",
			workingDirectory: "~/",
			wantErr:          false,
		},
		{
			testName:         "Root workingDirectory",
			name:             "Rooted",
			workingDirectory: "~/dotfiles/",
			wantErr:          false,
		},
		{
			// FIXME: this session does not close after the test
			testName:         "Non-existing workingDirectory",
			name:             "Non-existing",
			workingDirectory: "/i/do/not/exist",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tmux.NewSession(tt.name, tt.workingDirectory, true)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewSession() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewSession() succeeded unexpectedly")
			} else {
				if err := tmux.KillSession(tt.name); err != nil {
					t.Errorf("Error killing test session %s: %v", tt.name, err)
				}
			}
		})
	}
}

func TestHasSession(t *testing.T) {
	tests := []struct {
		testName string
		// Named input parameters for target function.
		name string
		want bool
	}{
		{
			testName: "Existing",
			name:     "Test",
			want:     true,
		},
		{
			testName: "Non-existing",
			name:     "I do not exist",
			want:     false,
		},
	}
	if err := tmux.NewSession("Test", "", true); err != nil {
		t.Errorf("Error opening test session: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tmux.HasSession(tt.name)
			if tt.want != got {
				t.Errorf("HasSession() = %v, want %v", got, tt.want)
			}
		})
	}
	if err := tmux.KillSession("Test"); err != nil {
		t.Errorf("Error killing test session: %v", err)
	}
}

func TestKillSession(t *testing.T) {
	tests := []struct {
		testName string
		// Named input parameters for target function.
		name    string
		wantErr bool
	}{
		{
			testName: "Kill existing",
			name:     "Test",
			wantErr:  false,
		},
		{
			testName: "Kill non-existing",
			name:     "I do not exist",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				if err := tmux.NewSession(tt.name, "", true); err != nil {
					t.Errorf("Error opening test session: %v", err)
				}
			}
			gotErr := tmux.KillSession(tt.name)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("KillSession() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("KillSession() succeeded unexpectedly")
			}
		})
	}
}

func TestSwitchToSession(t *testing.T) {
	tests := []struct {
		testName string
		// Named input parameters for target function.
		name    string
		wantErr bool
	}{
		{
			testName: "Existing",
			name:     "Test",
			wantErr:  false,
		},
		{
			testName: "Non-existing",
			name:     "I do not exist",
			wantErr:  true,
		},
	}
	currentSession := tmux.GetCurrentTmuxSessionName()
	if err := tmux.NewSession("Test", "", true); err != nil {
		t.Errorf("Error opening test session: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tmux.SwitchToSession(tt.name)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SwitchToSession() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SwitchToSession() succeeded unexpectedly")
			}
		})
	}
	if err := tmux.SwitchToSession(currentSession); err != nil {
		t.Errorf("Error switch back to tests execution session %s: %v", currentSession, err)
	}
	if err := tmux.KillSession("Test"); err != nil {
		t.Errorf("Error killing test session: %v", err)
	}
}

func TestGetCurrentTmuxSessionName(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Basic test",
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tmux.GetCurrentTmuxSessionName()
			if got != tt.want {
				t.Errorf("GetCurrentTmuxSessionName() = %v, want %v", got, tt.want)
			}
		})
	}
}
