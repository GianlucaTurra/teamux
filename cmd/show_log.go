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
	Long: `Show n log lines from teamux. Without "n" cat is invoked. 
	The user can clear the file with the clear flag (short -c)`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		clear, _ := cmd.Flags().GetBool("clear")
		if clear {
			err = common.ClearLogFile()
		} else {
			err = showLogs(args)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(showLogCmd)
	showLogCmd.Flags().BoolP("clear", "c", false, "Clear log file")
}

// showLogs runs either cat or tail on the log file and returns the combined
// output to the user
func showLogs(args []string) error {
	n := -1
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
}
