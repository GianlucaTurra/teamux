package cmd

import (
	"log"

	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/database"
	"github.com/spf13/cobra"
)

var newSessionCmd = &cobra.Command{
	Use:   "news [name]",
	Short: "Start and attach a new teamux session",
	Long: `Start a new tmux session and attach the tmux server to it. The session must be 
	present in the local teamux DB`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		detached, _ := cmd.Flags().GetBool("detached")
		if len(args) == 1 {
			name = args[0]
		}
		db, err := database.GetProdDB()
		if err != nil {
			log.Fatal(err)
		}
		session := sessions.Session{Name: name}
		db.Preload("Windows").Preload("Windows.Panes").First(&session)
		return session.Open(detached)
	},
}

func init() {
	newSessionCmd.Flags().BoolP("detached", "d", false, "Open the session as detached")
	rootCmd.AddCommand(newSessionCmd)
}
