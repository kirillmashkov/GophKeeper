package cmd

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// createCardCmd represents the createCard command
var createCardCmd = &cobra.Command{
	Use:   "card",
	Short: "create card secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal("error read secret name")
		}

		number, err := cmd.Flags().GetString("number")
		if err != nil {
			log.Fatal("error read number card for secret")
		}

		expiry, err := cmd.Flags().GetString("expiry")
		if err != nil {
			log.Fatal("error read expiry card for secret")
		}

		cvc, err := cmd.Flags().GetString("cvc")
		if err != nil {
			log.Fatal("error read cvc card for secret")
		}

		holder, err := cmd.Flags().GetString("holder")
		if err != nil {
			log.Fatal("error read holder card for secret")
		}

		metaDataFile, err := cmd.Flags().GetString("metadata")
		if err != nil {
			log.Fatal("error read metadata file for secret")
		}

		metaData, err := os.ReadFile(metaDataFile)
		if err != nil {
			log.Fatal("error read metadata file")
		}

		card := client.Card {
			Holder: holder,
			Cvc: cvc,
			Expiry: expiry,
			Number: number,
		}

			secret, err := client.Encode(card)
			if err != nil {
				log.Fatal("error create requet when create card secret")
			}

			dataEncoded := base64.StdEncoding.EncodeToString(secret)
			metaDataEncoded := base64.StdEncoding.EncodeToString(metaData)

			req := model.CreateSecretRequest {
				Name: name,
				Kind: card.Kind(),
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
	createCmd.AddCommand(createCardCmd)

	createCardCmd.Flags().StringP("name", "n", "", "Secret name")
	if err := createCardCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

	createCardCmd.Flags().StringP("number", "u", "", "Card number")
	if err := createCardCmd.MarkFlagRequired("number"); err != nil {
		log.Fatal(err)
	}

	createCardCmd.Flags().StringP("expiry", "e", "", "Card expiry date")
	if err := createCardCmd.MarkFlagRequired("expiry"); err != nil {
		log.Fatal(err)
	}

	createCardCmd.Flags().StringP("cvc", "v", "", "Card security code")
	if err := createCardCmd.MarkFlagRequired("cvc"); err != nil {
		log.Fatal(err)
	}

	createCardCmd.Flags().StringP("holder", "o", "", "Card holder")
	if err := createCardCmd.MarkFlagRequired("holder"); err != nil {
		log.Fatal(err)
	}

	createCardCmd.Flags().StringP("metadata", "m", "", "metadata file")
}
