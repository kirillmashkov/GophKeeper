package cmd

import (
	"log"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// createBinCmd represents the createBin command
var createBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "create bin secret",
	Run: func(cmd *cobra.Command, args []string) {
		name := mustGetStringFlag(cmd, "name", "error reaf secret name")
		file := mustGetStringFlag(cmd, "file", "error read file for secret")
		data := mustReadFile(file, "error read secret file")
		metaData := readMetadata(cmd)

		b := client.Bin{Data: data}
		createSecret(name, b, metaData, "bin")
	},
}

func init() {
	createCmd.AddCommand(createBinCmd)

	createBinCmd.Flags().StringP("name", "n", "", "Secret name")
	if err := createBinCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}
	createBinCmd.Flags().StringP("file", "f", "", "Binary file")
	if err := createBinCmd.MarkFlagRequired("file"); err != nil {
		log.Fatal(err)
	}

	createBinCmd.Flags().StringP("metadata", "m", "", "metadata file")
}
