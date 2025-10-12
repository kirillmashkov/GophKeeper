package cmd

import (
	"log"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"github.com/spf13/cobra"
)


var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login user",
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatal("Failed to read email")
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal("Failed to read password")
		}

		req := model.SignUpRequest {
			Email: email,
			Password: password,
		}

		r, err := keeperClient.R().SetBody(req).Post(addr + "/api/signin")
		if err != nil {
			log.Fatalf("error login user %s", email)
		}

		token := r.Header().Get("Authorization")
		err = tokenStorage.Save(token)
		if err != nil {
			log.Fatalf("error save token for user %s", email)
		}

		log.Printf("user was login %s", email)

	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("email", "e", "", "User Email")
	if err := loginCmd.MarkFlagRequired("email"); err != nil {
		log.Fatal(err)
	}
	loginCmd.Flags().StringP("password", "p", "", "User password")
	if err := loginCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err)
	}

}
