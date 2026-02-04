package cmd

import (
	"github.com/GianlucaTurra/teamux/database"
	"github.com/GianlucaTurra/teamux/tests"
	"github.com/spf13/cobra"
)

var testModeCmd = &cobra.Command{
	Use:   "test-mode",
	Short: "Launch with the test DB",
	Long:  `Launch teamux with a test a in memory DB that starts empty`,
	Args:  cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn := database.GetTestDB()
		err := tests.MigrateTestDB(conn.DB)
		if err != nil {
			return err
		}
		tui(conn.DB)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testModeCmd)
}
