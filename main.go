package main

import (
	"strconv"

	"github.com/aryanA101a/legoshichat-backend/handler"
	"github.com/aryanA101a/legoshichat-backend/store"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	authStore := store.NewAuthStore(app.DB())
	messageStore := store.NewMessageStore(app.DB())
	authCreator := handler.NewCreator()
	h := handler.Handler{Auth: authStore, Message: messageStore, AuthCreator: authCreator}

	app.POST("/create-account", h.HandleCreateAccount)
	app.POST("/login", h.HandleLogin)
	app.POST("/message/sendById",handler.WithJWTAuth(h.HandleSendMessageByID,authStore))
	app.POST("/message/sendByPhoneNumber",handler.WithJWTAuth(h.HandleSendMessageByPhoneNumber,authStore))

	port, err := strconv.Atoi(app.Config.Get("HTTP_PORT"))
	if err == nil {
		app.Server.HTTP.Port = port
	}

	app.Start()
}
