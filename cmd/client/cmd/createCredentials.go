package cmd

import (
	"log"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// createCredentialsCmd represents the createCredentials command
var createCredentialsCmd = &cobra.Command{
	Use:   "cred",
	Short: "create credentials secret",
	Run: func(cmd *cobra.Command, args []string) {
		name := mustGetStringFlag(cmd, "name", "error read secret name")
		login := mustGetStringFlag(cmd, "login", "error read login for secret")
		password := mustGetStringFlag(cmd, "password", "error read password for secret")

		metaData := readMetadata(cmd)

		cred := client.Credential{
			Login:    login,
			Password: password,
		}

		createSecret(name, cred, metaData, "credentials")
	},
}

func init() {
	createCmd.AddCommand(createCredentialsCmd)

	createCredentialsCmd.Flags().StringP("name", "n", "", "Secret name")
	if err := createCredentialsCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

	createCredentialsCmd.Flags().StringP("login", "l", "", "login")
	if err := createCredentialsCmd.MarkFlagRequired("login"); err != nil {
		log.Fatal(err)
	}

	createCredentialsCmd.Flags().StringP("password", "p", "", "password")
	if err := createCredentialsCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err)
	}

	createCredentialsCmd.Flags().StringP("metadata", "m", "", "metadata file")


}
