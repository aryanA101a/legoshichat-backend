package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreateAccountResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
type CreateAccountRequest struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber uint64 `json:"phoneNumber" validate:"required,min=1000000000,max=9999999999"`
	Password    string `json:"password" validate:"required,min=8"`
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
