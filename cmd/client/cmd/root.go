package cmd

import (
	"os"

	"github.com/kirillmashkov/GophKeeper.git/pkg/token"
	"github.com/spf13/cobra"
	"resty.dev/v3"
)

const addr = "http://localhost:8080"

var keeperClient *resty.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gophkeeper-cli",
	Short: "client",
}

const tokenFile = "token.txt"

var tokenStorage token.TokenStorage

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	tokenStorage = *token.NewTokenStorage(tokenFile)

	keeperClient = resty.New()

}


