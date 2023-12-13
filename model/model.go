package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	Name        string `json:"name"`
	PhoneNumber uint64 `json:"phoneNumber"`
	Password    string `json:"password"`
}
type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber uint64 `json:"phoneNumber"`
	Password    string `json:"password"`
}
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber uint64 `json:"phoneNumber"`
}

type Error struct {
	Message string
}

func NewAccount(name string, phoneNumber uint64, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:          uuid.New().String(),
		Name:        name,
		PhoneNumber: phoneNumber,
		Password:    string(encpw),
	}, nil
}
