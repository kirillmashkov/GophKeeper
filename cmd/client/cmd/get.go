package cmd

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal("error reaf secret name")
		}

		var resp model.GetSecretResponse
		req := fmt.Sprintf("%s/api/%s", addr, name)
		r, err := keeperClient.R().SetResult(&resp).SetHeader("Authorization", tokenStr).Get(req)

		if err != nil {
			log.Fatalf("error get secret %s", name)
		}

		if r.StatusCode() != 200 {
			log.Fatalf("error get secret %s http code %d", name, r.StatusCode())
		}

		dataByte, err := base64.StdEncoding.DecodeString(resp.Data)
		if err != nil {
			log.Fatalf("error decode secret %s", name)
		}

		s, err := client.Decode(dataByte)
		if err != nil {
			log.Fatalf("error parse secret %s", name)
		}

		log.Printf("%s", s.String())

	},
}

func init() {
	secretCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("name", "n", "", "Secret name")
	if err := getCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}

}
