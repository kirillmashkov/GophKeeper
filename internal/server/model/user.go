package model

import "errors"

// User содержит УЗ пользователя
type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type SignUpRequest struct {
	Email        string `json:"email"`
	Password 	 string `json:"password"`
}

// Возможные ошибки при работе с хранилищем UserStorage
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserConflict = errors.New("user conflict")
	ErrSecretConflict = errors.New("secret conflict")
	ErrSecretNoFound = errors.New("secret not found")
)