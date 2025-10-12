package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const key = "keyforhash"

func Hash(data string) (string, error) {
	mac := hmac.New(sha256.New, []byte(key))
	if _, err := mac.Write([]byte(data)); err != nil {
		return "", err
	}
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum), nil
}
