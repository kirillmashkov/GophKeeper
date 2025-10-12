package cmd

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// createCredentialsCmd represents the createCredentials command
var createCredentialsCmd = &cobra.Command{
	Use:   "cred",
	Short: "create credentials secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal("error read secret name")
		}

		login, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Fatal("error read login for secret")
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal("error read password for secret")
		}

		metaDataFile, err := cmd.Flags().GetString("metadata")
		if err != nil {
			log.Fatal("error read metadata file for secret")
		}

		metaData, err := os.ReadFile(metaDataFile)
		if err != nil {
			log.Fatal("error read metadata file")
		}

		cred := client.Credential {
			Login: login,
			Password: password,
		}

		secret, err := client.Encode(cred)
		if err != nil {
			log.Fatal("error create request when create credentials secret")
		}

		dataEncoded := base64.StdEncoding.EncodeToString(secret)
		metaDataEncoded := base64.StdEncoding.EncodeToString(metaData)

		req := model.CreateSecretRequest {
			Name: name,
			Kind: cred.Kind(),
			Data: dataEncoded,
			Metadata: metaDataEncoded,
		}

		r, err := keeperClient.R().SetHeader("Authorization", tokenStr).SetBody(req).Post(addr + "/api/createsecret")

		if err != nil {
			log.Fatalf("error create card secret %s", name)
		}

		if r.StatusCode() == 409 {
			log.Printf("card secret %s already exist", name)	
		} else if r.StatusCode() == 201 {
			log.Printf("create card secret %s", name)
		} else {
			log.Fatalf("error create card secret %s", name)
		}

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
