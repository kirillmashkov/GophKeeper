package cmd

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// createBinCmd represents the createBin command
var createBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "create bin secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal("error reaf secret name")
		}

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal("error read file for secret")
		}

		metaDataFile, err := cmd.Flags().GetString("metadata")
		if err != nil {
			log.Fatal("error read file for secret")
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal("error read secret file")
		}

		metaData, err := os.ReadFile(metaDataFile)
		if err != nil {
			log.Fatal("error read metadata file")
		}

		b := client.Bin{Data: data}
		secret, err := client.Encode(b)
		if err != nil {
			log.Fatal("error create requet when create bin secret")
		}

		dataEncoded := base64.StdEncoding.EncodeToString(secret)
		metaDataEncoded := base64.StdEncoding.EncodeToString(metaData)

		req := model.CreateSecretRequest {
			Name: name,
			Kind: b.Kind(),
			Data: dataEncoded,
			Metadata: metaDataEncoded,
		}

		r, err := keeperClient.R().SetHeader("Authorization", tokenStr).SetBody(req).Post(addr + "/api/createsecret")

		if err != nil {
			log.Fatalf("error create bin secret %s", name)
		}

		if r.StatusCode() == 409 {
			log.Printf("bin secret %s already exist", name)	
		} else if r.StatusCode() == 201 {
			log.Printf("create bin secret %s", name)
		} else {
			log.Fatalf("error create bin secret %s", name)
		}
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
