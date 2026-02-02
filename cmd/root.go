// cmd contains all the CLI commands available in teamux.
// Commands are defined using cobra cli
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components"
	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const LogFile = "/tmp/teamux.log"

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Easy tmux sessions",
	Long: `Teamux is a TUI application to define tmux sessions, windows and 
	panes and manage them.`,
	Run: func(cmd *cobra.Command, args []string) { tui() },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

// tui open the TUI for teamux, creating the temporary log file and ensuring
// db migrations are applied.
func tui() {
	setup()
	logfile, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalln("Unable to setup syslog:", err.Error())
	}
	defer logfile.Close()
	teamuxLogger := common.Logger{
		Infologger:    log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warninglogger: log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Errorlogger:   log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Fatallogger:   log.New(logfile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	db, err := database.GetProdDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(
		&sessions.Session{},
		&windows.Window{},
		&panes.Pane{},
	); err != nil {
		teamuxLogger.Errorlogger.Printf("Error migrating tables: %v", err)
	}
	p := tea.NewProgram(components.InitialModel(database.Connector{DB: db, Ctx: context.Background()}, teamuxLogger))
	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

// setup checks if the config directory exists, if it doesn't it creates it
func setup() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error detecting user config directory: %v", err)
	}
	teamuxConfigDir := fmt.Sprintf("%s/%s", userConfigDir, "teamux")
	if _, err := os.Stat(teamuxConfigDir); err == nil {
		return
	}
	if err := os.Mkdir(teamuxConfigDir, 0o755); err != nil {
		log.Fatalf("Error creating config dir: %v", err)
	}
}
