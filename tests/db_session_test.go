package tests

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/GianlucaTurra/teamux/components/data"
	"gorm.io/gorm"
)

func sampleSession(name string, connector data.Connector) error {
	_, err := data.CreateSession(name, "", connector)
	return err
}

func TestOpenSession(t *testing.T) {
	connector := setup()
	name := "OpenTest"
	closeTestSession(name)
	// Count current tmux sessions, this is needed since the tmux server might be already running
	n, err := countTmuxSession()
	if err != nil {
		t.Fatalf("Failed to count tmux sessions: %v", err)
	}
	// Create new session and open it
	if err := sampleSession("OpenTest", connector); err != nil {
		t.Errorf("Failed to create sample session: %v", err)
	}
	s, err := gorm.G[data.Session](connector.DB).Where("id = ?", 1).First(connector.Ctx)
	if err != nil {
		t.Errorf("Failed to read session: %v", err)
	}
	if err := s.Open(); err != nil {
		t.Errorf("Failed to open session: %v", err)
	}
	// Compare new number of open sessions
	m, err := countTmuxSession()
	if err != nil {
		t.Fatalf("Failed to count tmux sessions: %v", err)
	}
	if m != n+1 {
		t.Errorf("Expected %d tmux sessions, got %d", n+1, m)
	}
	closeTestSession(name)
}

// createSampleWindows Creates 3 sample windows in the expiring db
func createSampleWindows(connector data.Connector) error {
	for i := range 3 {
		if _, err := data.CreateWindow(fmt.Sprintf("Window %d", i), "", connector); err != nil {
			return err
		}
	}
	return nil
}

// TODO: might be useful in the entire app
func closeTestSession(name string) error {
	checkCmd := exec.Command("sh", "-c", fmt.Sprintf("tmux has-session -t %s", name))
	if err := checkCmd.Run(); err != nil {
		return nil
	}
	closeCmd := exec.Command("sh", "-c", fmt.Sprintf("tmux kill-session -t %s", name))
	return closeCmd.Run()
}

// TODO: might be useful in the entire app
func countTmuxSession() (int, error) {
	cmd := exec.Command("sh", "-c", "tmux ls | wc -l")
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to count tmux sessions: %w", err)
	}
	var n int
	if _, err := fmt.Sscanf(string(out), "%d", &n); err != nil {
		return 0, fmt.Errorf("failed to parse tmux session count: %w", err)
	}
	return n, nil
}
