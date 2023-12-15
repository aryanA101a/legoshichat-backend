package handler

import (
	"database/sql"
	"encoding/json"
	"time"

	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/aryanA101a/legoshichat-backend/store"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/types"
	"golang.org/x/crypto/bcrypt"
)

type Creator interface {
	NewAccount(name string, phoneNumber uint64, password string) (*model.Account, error)
	NewJWT(ctx *gofr.Context, userId string) (string, error)
}

type authCreator struct{}

type handler struct {
	auth        store.AuthStore
	authCreator Creator
}

func New(a store.AuthStore,c Creator) handler {
	return handler{auth: a, authCreator: c}
}

func NewCreator()authCreator{
	return authCreator{}
}

func (h handler) HandleCreateAccount(ctx *gofr.Context) (interface{}, error) {
	var accountRequest model.CreateAccountRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&accountRequest)
	err = validator.New().Struct(accountRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields"+err.Error())
	}

	newAccount, err := h.authCreator.NewAccount(accountRequest.Name, accountRequest.PhoneNumber, accountRequest.Password)
	if err != nil {
		return nil, e.HttpStatusError(500, "")
	}

	user, err := h.auth.CreateAccount(ctx, *newAccount)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, err.Error())
	}

	token, err := h.authCreator.NewJWT(ctx, user.ID)
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
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields -"+err.Error())
	}

	account, err := h.auth.FetchAccount(ctx, loginRequest)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(401, "User does not exists")
		}
		return nil, e.HttpStatusError(500, "")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginRequest.Password))
	if err != nil {
		return nil, e.HttpStatusError(401, "Invalid password")
	}

	token, err := h.authCreator.NewJWT(ctx, account.ID)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, "")
	}
	return types.Raw{Data: model.AuthResponse{User: model.User{ID: account.ID, Name: account.Name, PhoneNumber: account.PhoneNumber}, Token: token}}, nil
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
		"expiresAt": time.Now().Add(time.Duration(time.Duration.Hours(24) * 5)),
	}

	secret := ctx.Config.Get("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
