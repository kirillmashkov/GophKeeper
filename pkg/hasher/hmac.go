package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

const key = "keyforhash"

func Hash(data string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
