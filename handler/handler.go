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
	Friend      store.FriendStore
	AuthCreator Creator
}

func New(a store.AuthStore, m store.MessageStore, f store.FriendStore, c Creator) Handler {
	return Handler{Auth: a, Message: m, Friend: f, AuthCreator: c}
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

	h.Friend.AddFriend(ctx, ctx.Value("userId").(string), messageRequest.RecipientID)

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

	recipientId, err := h.Auth.GetUserIdByPhoneNumber(ctx, messageRequest.RecipientPhoneNumber)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(404, "Recipient does not exists")
		}
		return nil, e.HttpStatusError(500, "")
	}

	message := NewMessage(ctx.Value("userId").(string), *recipientId, messageRequest.Content)

	h.Friend.AddFriend(ctx, ctx.Value("userId").(string), *recipientId)

	err = h.Message.AddMessage(ctx, message)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(500, err.Error())
	}

	return types.Raw{Data: message}, nil
}

func (h Handler) HandleGetMessage(ctx *gofr.Context) (interface{}, error) {
	messageId := ctx.PathParam("id")
	if strings.TrimSpace(messageId) == "" {
		return nil, e.HttpStatusError(400, "Missing Parameter messageId")
	}

	message, err := h.Message.GetMessage(ctx, ctx.Value("userId").(string), messageId)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(404, "Message does not exists")
		} else if err == e.NewError("You are not authorized to see that message") {
			return nil, e.HttpStatusError(403, err.Error())
		}
		return nil, e.HttpStatusError(500, "")
	}
	return types.Raw{Data: message}, nil
}

func (h Handler) HandlePutMessage(ctx *gofr.Context) (interface{}, error) {
	var updateMessageRequest model.UpdateMessageRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&updateMessageRequest)
	err = validator.New().Struct(updateMessageRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields -"+err.Error())
	}

	messageId := ctx.PathParam("id")
	if strings.TrimSpace(messageId) == "" {
		return nil, e.HttpStatusError(400, "Missing Parameter messageId")
	}

	message, err := h.Message.UpdateMessage(ctx, ctx.Value("userId").(string), messageId, updateMessageRequest.Content)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(404, "Message does not exists")
		} else if err == e.NewError("You are not authorized to update that message") {
			return nil, e.HttpStatusError(403, err.Error())
		}
		return nil, e.HttpStatusError(500, "")
	}
	return types.Raw{Data: message}, nil
}

func (h Handler) HandleDeleteMessage(ctx *gofr.Context) (interface{}, error) {
	messageId := ctx.PathParam("id")
	if strings.TrimSpace(messageId) == "" {
		return nil, e.HttpStatusError(400, "Missing Parameter messageId")
	}

	err := h.Message.DeleteMessage(ctx, ctx.Value("userId").(string), messageId)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(404, "Message does not exists")
		} else if err == e.NewError("You are not authorized to delete that message") {
			return nil, e.HttpStatusError(403, err.Error())
		}
		return nil, e.HttpStatusError(500, "")
	}
	return nil, nil
}

func (h Handler) HandleGetMessages(ctx *gofr.Context) (interface{}, error) {
	var getMessageRequest model.GetMessagesRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&getMessageRequest)
	err = validator.New().Struct(getMessageRequest)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, e.HttpStatusError(400, "Invalid inputs or missing required fields -"+err.Error())
	}

	if userId := ctx.Value("userId").(string); !(userId == getMessageRequest.SenderID || userId == getMessageRequest.RecipientID) {
		return nil, e.HttpStatusError(403, "You are not authorized to retrieve these messages")
	}

	messages, err := h.Message.GetMessages(ctx, getMessageRequest.SenderID, getMessageRequest.RecipientID, getMessageRequest.Page, model.RequestMessageLimit)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(404, "No messages found for this page")
		}
		return nil, e.HttpStatusError(500, "")
	}

	if len(*messages) == 0 {
		return nil, e.HttpStatusError(404, "Page does not exists")
	}

	lastPage := false
	if len(*messages) < model.RequestMessageLimit {
		lastPage = true
	}
	return types.Raw{Data: model.GetMessagesResponse{Page: getMessageRequest.Page, LastPage: lastPage, Messages: *messages}}, nil

}

func (h Handler) GetFriends(ctx *gofr.Context) (interface{}, error) {
	friends, err := h.Friend.GetFriends(ctx, ctx.Value("userId").(string))
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil, e.HttpStatusError(404, "No Friends")
		}
		return nil, e.HttpStatusError(500, "")
	}
	return types.Raw{Data: friends}, nil
}

func WithJWTAuth(handlerFunc gofr.Handler, authStore store.AuthStore) gofr.Handler {
	return func(ctx *gofr.Context) (interface{}, error) {

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
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header missing")
	}

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
