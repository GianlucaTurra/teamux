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
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Easy tmux sessions",
	Long: `Teamux is a TUI application to define tmux sessions, windows and 
	panes and manage them.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.GetProdDB()
		if err != nil {
			log.Fatal(err)
		}
		if err = db.AutoMigrate(
			&sessions.Session{},
			&windows.Window{},
			&panes.Pane{},
		); err != nil {
			log.Fatal(err)
		}
		tui(db)
	},
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
func tui(db *gorm.DB) {
	setup()
	teamuxLogger := common.GetLogger()
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
