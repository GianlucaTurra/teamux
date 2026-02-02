package cmd

import (
	"log"

	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/database"
	"github.com/spf13/cobra"
)

var newPaneCmd = &cobra.Command{
	Use:   "newp [name]",
	Short: "Create a new pane for the current window",
	Long:  `Create a new pane for the current tmux window, using the current pane as a target. The pane must be present in the temaux DB`,
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
		pane := panes.Pane{Name: name}
		db.First(&pane)
		return pane.Open()
	},
}

func init() {
	rootCmd.AddCommand(newPaneCmd)
}
