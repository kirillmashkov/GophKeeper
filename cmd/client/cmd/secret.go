package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var tokenStr string

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secret",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		tokenStr, err = tokenStorage.Load()
		if err != nil {
			log.Fatal("error read token file")
		}
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
}
