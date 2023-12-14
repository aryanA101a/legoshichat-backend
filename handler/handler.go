package handler

import (
	"encoding/json"
	"time"

	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/aryanA101a/legoshichat-backend/store"
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

func (h handler) HandleCreateAccount(ctx *gofr.Context) (interface{}, error) {
	var accountRequest model.CreateAccountRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&accountRequest)
	err = validator.New().Struct(accountRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields"+err.Error())
	}

	newAccount, err := model.NewAccount(accountRequest.Name, accountRequest.PhoneNumber, accountRequest.Password)
	if err != nil {
		return nil, e.HttpStatusError(500, "")
	}

	user, err := h.auth.CreateAccount(ctx, *newAccount)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, err.Error())
	}

	token, err := createJWT(ctx, user.ID)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, "")
	}

	return types.Raw{Data: model.AuthResponse{User: *user, Token: token}}, nil
}

func (h handler) HandleLogin(ctx *gofr.Context) (interface{}, error) {
	var loginRequest model.LoginRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&loginRequest)
	err = validator.New().Struct(loginRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields -"+ err.Error())
	}

	user, err := h.auth.Login(ctx, loginRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, err
	}

	token, err := createJWT(ctx, user.ID)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, "")
	}
	return types.Raw{Data: model.AuthResponse{User: *user, Token: token}}, nil
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
