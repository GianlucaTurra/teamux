package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var NewSessionCmd = &cobra.Command{
	Use:   "news [name]",
	Short: "Start and attach a new teamux session",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "ASD"
		if len(args) == 1 {
			name = args[0]
		}
		fmt.Printf("HI %s", name)
		return nil
	},
}
