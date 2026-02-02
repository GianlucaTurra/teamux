package cmd

import (
	"log"

	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/database"
	"github.com/spf13/cobra"
)

var newWindowCmd = &cobra.Command{
	Use:   "neww [name]",
	Short: "Create a new window for the current session",
	Long:  `Create a new window for the current tmux session. The window must be present in the temaux DB`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) == 1 {
			name = args[0]
		}
		db, err := database.GetProdDB()
		if err != nil {
			log.Fatal(err)
		}
		window := windows.Window{Name: name}
		db.Preload("Panes").First(&window)
		return window.Open()
	},
}

func init() {
	rootCmd.AddCommand(newWindowCmd)
}
