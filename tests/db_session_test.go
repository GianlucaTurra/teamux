package tests

import (
	"database/sql"
	"fmt"
	"os/exec"
	"testing"

	"github.com/GianlucaTurra/teamux/internal/data"
	_ "github.com/mattn/go-sqlite3"
)

func sampleSession(name string, db *sql.DB) error {
	return data.NewSession(name, "", db).Save()
}

func TestSessionCreation(t *testing.T) {
	db := setup()
	defer db.Close()
	s := data.NewSession("Test", "", db)
	if err := s.Save(); err != nil {
		t.Errorf("Failed to save session: %v", err)
	}
	query := "SELECT COUNT(*) FROM sessions"
	row := db.QueryRow(query)
	var n int
	if err := row.Scan(&n); err != nil {
		t.Errorf("Failed to query session count: %v", err)
	}
	if n != 1 {
		t.Errorf("Expected 1 session, got %d", n)
	}
}

func TestSessionUpdate(t *testing.T) {
	db := setup()
	defer db.Close()
	if err := sampleSession("test", db); err != nil {
		t.Errorf("Failed to create sample session: %v", err)
	}
	s, err := data.ReadSessionByID(db, 1)
	if err != nil {
		t.Errorf("Failed to read session: %v", err)
	}
	s.Name = "UpdatedTest"
	if err := s.Save(); err != nil {
		t.Errorf("Failed to update session: %v", err)
	}
	s, err = data.ReadSessionByID(db, 1)
	if err != nil {
		t.Errorf("Failed to read session: %v", err)
	}
	if s.Name != "UpdatedTest" {
		t.Errorf("Expected session name 'UpdatedTest', got '%s'", s.Name)
	}
}

func TestReadAllSessions(t *testing.T) {
	db := setup()
	defer db.Close()
	for i := range 5 {
		if err := sampleSession(fmt.Sprintf("Session %d", i), db); err != nil {
			t.Errorf("Failed to create sample session: %v", err)
		}
	}
	sessions, err := data.ReadAllSessions(db)
	if err != nil {
		t.Errorf("Failed to read all sessions: %v", err)
	}
	if sessions == nil || len(sessions) != 5 {
		t.Errorf("Expected 5 sessions, got %d", len(sessions))
	}
}

func TestDeleteSession(t *testing.T) {
	db := setup()
	defer db.Close()
	if err := sampleSession("ToDelete", db); err != nil {
		t.Errorf("Failed to create sample session: %v", err)
	}
	s, err := data.ReadSessionByID(db, 1)
	if err != nil {
		t.Errorf("Failed to read session: %v", err)
	}
	if err := s.Delete(); err != nil {
		t.Errorf("Failed to delete session: %v", err)
	}
	_, err = data.ReadSessionByID(db, 1)
	if err == nil {
		t.Error("Expected error when reading deleted session, got none")
	}
}

func TestOpenSession(t *testing.T) {
	db := setup()
	defer db.Close()
	name := "OpenTest"
	closeTestSession(name)
	// Count current tmux sessions, this is needed since the tmux server might be already running
	n, err := countTmuxSession()
	if err != nil {
		t.Fatalf("Failed to count tmux sessions: %v", err)
	}
	// Create new session and open it
	if err := sampleSession("OpenTest", db); err != nil {
		t.Errorf("Failed to create sample session: %v", err)
	}
	s, err := data.ReadSessionByID(db, 1)
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

func TestGetAllWindows(t *testing.T) {
	db := setup()
	defer db.Close()
	name := "OpenTest"
	closeTestSession(name)
	// Create sample Session and open it
	if err := sampleSession(name, db); err != nil {
		t.Fatalf("Failed to create sample sessions: %v", err)
	}
	s, err := data.ReadSessionByID(db, 1)
	if err != nil {
		t.Errorf("Failed to read session: %v", err)
	}
	if err := s.Open(); err != nil {
		t.Errorf("Failed to open session: %v", err)
	}
	// Create sample windows and associate them to the session in db
	if err := createSampleWindows(db); err != nil {
		t.Errorf("Failed to create sample windows: %v", err)
	}
	if err := associateWindowsToSession(db); err != nil {
		t.Errorf("Failed to associate windows to session: %v", err)
	}
	windowIds, err := s.GetAllWindows()
	if err != nil {
		t.Errorf("Failed to read related windows: %v", err)
	}
	if len(windowIds) != 3 {
		t.Errorf("Expected 3 windows found: %d", len(windowIds))
	}
}

// createSampleWindows Creates 3 sample windows in the expiring db
func createSampleWindows(db *sql.DB) error {
	for i := range 3 {
		if err := data.NewWindow(fmt.Sprintf("Window %d", i), "", db).Save(); err != nil {
			return err
		}
	}
	return nil
}

func associateWindowsToSession(db *sql.DB) error {
	query := "INSERT INTO Session_Windows (session_id, window_id) VALUES (?, ?)"
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := range 3 {
		_, err := stmt.Exec(1, i)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
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
