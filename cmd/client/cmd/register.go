/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register User",
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatal("error email")
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatal("error password")
		}

		req := model.SignUpRequest{
			Email:    email,
			Password: password,
		}
		_, err = keeperClient.R().SetBody(req).Post(addr + "/api/signup")
		if err != nil {
			log.Fatalf("error create user %s", email)
		}

		log.Printf("user was created %s", email)
	},
}

func init() {
	authCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("email", "e", "", "User Email")
	if err := registerCmd.MarkFlagRequired("email"); err != nil {
		log.Fatal(err)
	}

	registerCmd.Flags().StringP("password", "p", "", "User password")
	if err := registerCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err)
	}
}
