package handler

import (
	"encoding/json"
	"time"

	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/aryanA101a/legoshichat-backend/store"
	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/types"
)

type handler struct {
	auth store.AuthStore
}

func New(a store.AuthStore) handler {
	return handler{auth: a}
}

func (h handler) CreateAccount(ctx *gofr.Context) (interface{}, error) {
	var accountRequest model.CreateAccountRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&accountRequest)
	err = validator.New().Struct(accountRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields")
	}

	ctx.Logger.Info(accountRequest)

	user, err := h.auth.CreateAccount(ctx, accountRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, err.Error())
	}
	token, err := createJWT(ctx, user.ID)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, "")
	}
	return types.Raw{Data: model.CreateAccountResponse{User: *user, Token: token}}, nil
}

func createJWT(ctx *gofr.Context, userId string) (string, error) {
	claims := &jwt.MapClaims{
		"id":        userId,
		"expiresAt": time.Now().Add(time.Duration(time.Duration.Hours(24) * 5)),
	}

	secret := ctx.Config.Get("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
