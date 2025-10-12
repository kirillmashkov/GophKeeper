package token

import (
	"io"
	"log"
	"os"
)

type TokenStorage struct {
	Path string
}

func NewTokenStorage(path string) *TokenStorage {
	return &TokenStorage{
		Path:  path,
	}
}

func (ts *TokenStorage) Save(token string) error {
	file, err := os.Create(ts.Path)
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Print("error close file with token")
		}
	}()

	_, err = file.WriteString(token)
	return err
}

func (ts *TokenStorage) Load() (string, error) {
	file, err := os.Open(ts.Path)
	if err != nil {
		return "", nil
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Print("error close file with token")
		}
	}()

	b, err := io.ReadAll(file)
	return string(b), err
}