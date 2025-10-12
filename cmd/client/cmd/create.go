package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create secret",
}

func init() {
	secretCmd.AddCommand(createCmd)
}
