package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/aryanA101a/legoshichat-backend/store"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/types"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Auth        store.AuthStore
	Message     store.MessageStore
	AuthCreator Creator
}

func New(a store.AuthStore, m store.MessageStore, c Creator) Handler {
	return Handler{Auth: a, Message: m, AuthCreator: c}
}

func (h Handler) HandleCreateAccount(ctx *gofr.Context) (interface{}, error) {
	var accountRequest model.CreateAccountRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&accountRequest)
	err = validator.New().Struct(accountRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields"+err.Error())
	}

	newAccount, err := h.AuthCreator.NewAccount(accountRequest.Name, accountRequest.PhoneNumber, accountRequest.Password)
	if err != nil {
		return nil, e.HttpStatusError(500, "")
	}

	user, err := h.Auth.CreateAccount(ctx, *newAccount)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, err.Error())
	}

	token, err := h.AuthCreator.NewJWT(ctx, user.ID)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, "")
	}

	return types.Raw{Data: model.AuthResponse{User: *user, Token: token}}, nil
}

func (h Handler) HandleLogin(ctx *gofr.Context) (interface{}, error) {
	var loginRequest model.LoginRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&loginRequest)
	err = validator.New().Struct(loginRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields -"+err.Error())
	}

	account, err := h.Auth.FetchAccount(ctx, loginRequest)
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

	token, err := h.AuthCreator.NewJWT(ctx, account.ID)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, "")
	}
	return types.Raw{Data: model.AuthResponse{User: model.User{ID: account.ID, Name: account.Name, PhoneNumber: account.PhoneNumber}, Token: token}}, nil
}

func (h Handler) HandleSendMessageByID(ctx *gofr.Context) (interface{}, error) {

	var messageRequest model.SendMessageByIDRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&messageRequest)
	err = validator.New().Struct(messageRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields -"+err.Error())
	}

	message := NewMessage(ctx.Value("userId").(string), messageRequest.RecipientID, messageRequest.Content)

	err = h.Message.AddMessage(ctx, message)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, err.Error())
	}

	return types.Raw{Data: message}, nil
}

func (h Handler) HandleSendMessageByPhoneNumber(ctx *gofr.Context) (interface{}, error) {

	var messageRequest model.SendMessageByPhoneNumberRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&messageRequest)
	err = validator.New().Struct(messageRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields -"+err.Error())
	}
	recipientId,err:=h.Auth.GetUserIdByPhoneNumber(ctx,messageRequest.RecipientPhoneNumber)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(401, "Recipient does not exists")
		}
		return nil, e.HttpStatusError(500, "")
	}
	message := NewMessage(ctx.Value("userId").(string), *recipientId, messageRequest.Content)

	err = h.Message.AddMessage(ctx, message)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, err.Error())
	}

	return types.Raw{Data: message}, nil
}

func WithJWTAuth(handlerFunc gofr.Handler, authStore store.AuthStore) gofr.Handler {
	return func(ctx *gofr.Context) (interface{}, error) {
		fmt.Println("calling JWT auth middleware")

		tokenString, err := extractToken(ctx)
		if err != nil {
			return nil, e.HttpStatusError(401, "Invalid token format")
		}
		token, err := parseJWT(ctx, tokenString)
		if err != nil {
			return nil, e.HttpStatusError(401, "Invalid token format")
		}
		if !token.Valid {
			return nil, e.HttpStatusError(401, "JWT token not valid")
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := claims["id"].(string)

		ok, err := authStore.AccountExists(ctx, userID)
		if err != nil {
			return nil, e.HttpStatusError(500, "")
		}
		if !ok {
			return nil, e.HttpStatusError(401, "User does not exist")
		}

		*&ctx.Context = context.WithValue(ctx.Context, "userId", userID)
		return handlerFunc(ctx)
	}
}

func extractToken(ctx *gofr.Context) (string, error) {
	// Get the Authorization header
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header missing")
	}

	// The header should be in the format "Bearer <token>"
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return "", fmt.Errorf("Invalid Authorization header format")
	}

	return splitToken[1], nil
}

func parseJWT(ctx *gofr.Context, tokenString string) (*jwt.Token, error) {
	secret := ctx.Config.Get("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
