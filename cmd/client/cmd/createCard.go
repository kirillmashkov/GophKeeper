package cmd

import (
	"log"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// createCardCmd represents the createCard command
var createCardCmd = &cobra.Command{
	Use:   "card",
	Short: "create card secret",
	Run: func(cmd *cobra.Command, args []string) {
		name := mustGetStringFlag(cmd, "name", "error read secret name")
		number := mustGetStringFlag(cmd, "number", "error read number card for secret")
		expiry := mustGetStringFlag(cmd, "expiry", "error read expiry card for secret")
		cvc := mustGetStringFlag(cmd, "cvc", "error read cvc card for secret")
		holder := mustGetStringFlag(cmd, "holder", "error read holder card for secret")

		metaData := readMetadata(cmd)

		card := client.Card{
			Holder: holder,
			Cvc:    cvc,
			Expiry: expiry,
			Number: number,
		}

		createSecret(name, card, metaData, "card")
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
