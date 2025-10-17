package cmd

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model/client"
	"github.com/spf13/cobra"
)

// mustGetStringFlag reads a string flag or terminates the process with a descriptive error.
func mustGetStringFlag(cmd *cobra.Command, flagName, errMsg string) string {
	v, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Fatal(errMsg)
	}
	return v
}

// mustReadFile reads a file or terminates the process with a descriptive error.
func mustReadFile(path, errMsg string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(errMsg)
	}
	return b
}

// readMetadata reads the metadata file path from the "metadata" flag and returns its contents.
func readMetadata(cmd *cobra.Command) []byte {
	metaFile := mustGetStringFlag(cmd, "metadata", "error read metadata file for secret")
	return mustReadFile(metaFile, "error read metadata file")
}

// createSecret encodes and sends a create-secret request and logs the result consistently.
func createSecret(name string, s client.ISecret, metadata []byte, secretType string) {
	secret, err := client.Encode(s)
	if err != nil {
		log.Fatalf("error create request when create %s secret", secretType)
	}

	dataEncoded := base64.StdEncoding.EncodeToString(secret)
	metaDataEncoded := base64.StdEncoding.EncodeToString(metadata)

	req := model.CreateSecretRequest{
		Name:     name,
		Kind:     s.Kind(),
		Data:     dataEncoded,
		Metadata: metaDataEncoded,
	}

	r, err := keeperClient.R().SetHeader("Authorization", tokenStr).SetBody(req).Post(addr + "/api/createsecret")
	if err != nil {
		log.Fatalf("error create %s secret %s", secretType, name)
	}

	switch r.StatusCode() {
	case 409:
		log.Printf("%s secret %s already exist", secretType, name)
	case 201:
		log.Printf("create %s secret %s", secretType, name)
	default:
		log.Fatalf("error create %s secret %s", secretType, name)
	}
}
