package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type Creator interface {
	NewAccount(name string, phoneNumber uint64, password string) (*model.Account, error)
	NewJWT(ctx *gofr.Context, userId string) (string, error)
}

type authCreator struct{}

func NewCreator() authCreator {
	return authCreator{}
}

func (authCreator) NewAccount(name string, phoneNumber uint64, password string) (*model.Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &model.Account{
		ID:          uuid.New().String(),
		Name:        name,
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
	}, nil
}

func (authCreator) NewJWT(ctx *gofr.Context, userId string) (string, error) {
	claims := &jwt.MapClaims{
		"id":        userId,
		"expiresAt": time.Now().Add(time.Duration(time.Hour * 24)),
	}

	secret := ctx.Config.Get("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func NewMessage(from, to, content string) model.Message {

	return model.Message{
		ID:      uuid.New().String(),
		From:    from,
		To:      to,
		Content: content,
		Timestamp: time.Now(),
	}
}
