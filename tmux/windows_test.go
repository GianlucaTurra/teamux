package tmux_test

import (
	"testing"

	"github.com/GianlucaTurra/teamux/tmux"
)

func TestNewWindow(t *testing.T) {
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
			testName:         "Non-existing workingDirectory",
			name:             "Non-existing",
			workingDirectory: "/i/do/not/exist",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tmux.NewWindow(tt.name, tt.workingDirectory)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateWindow() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateWindow() succeeded unexpectedly")
			}
			if err := tmux.KillWindow(tt.name); err != nil {
				t.Errorf("Error killing test session: %v", err)
			}
		})
	}
}

func TestNewWindowWithTarget(t *testing.T) {
	tests := []struct {
		testName string
		// Named input parameters for target function.
		name             string
		workingDirectory string
		target           string
		wantErr          bool
	}{
		{
			testName:         "Empty workingDirectory",
			name:             "Test",
			workingDirectory: "",
			target:           "Test",
			wantErr:          false,
		},
		{
			testName:         "Home workingDirectory",
			name:             "Home",
			workingDirectory: "~/",
			target:           "Test",
			wantErr:          false,
		},
		{
			testName:         "Root workingDirectory",
			name:             "Rooted",
			workingDirectory: "~/dotfiles/",
			target:           "Test",
			wantErr:          false,
		},
		{
			testName:         "Non-existing workingDirectory",
			name:             "Non-existing",
			workingDirectory: "/i/do/not/exist",
			target:           "Test",
			wantErr:          true,
		},
	}
	if err := tmux.NewSession("Test", ""); err != nil {
		t.Errorf("Error opening test session: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tmux.NewWindowWithTarget(tt.name, tt.workingDirectory, tt.target)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateWindowWithTarget() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateWindowWithTarget() succeeded unexpectedly")
			}
		})
	}
	if err := tmux.KillSession("Test"); err != nil {
		t.Errorf("Error killing test session: %v", err)
	}
}
