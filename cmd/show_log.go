package cmd

import (
	"fmt"
	"strconv"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/spf13/cobra"
)

var showLogCmd = &cobra.Command{
	Use:   "logs [n]",
	Short: "Show n log lines",
	Long:  "Show n log lines from teamux",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 5
		var err error
		if len(args) == 1 {
			if n, err = strconv.Atoi(args[0]); err != nil {
				return err
			}
		}
		if bytes, err := common.ShowLogFile(n); err != nil {
			return err
		} else {
			fmt.Print(string(bytes))
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(showLogCmd)
}
