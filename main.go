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
    authCreator:=handler.NewCreator()
	handler := handler.New(authStore,authCreator)
	
    
    app.POST("/create-account", handler.HandleCreateAccount)
    app.POST("/login", handler.HandleLogin)
    
    port,err:=strconv.Atoi(app.Config.Get("HTTP_PORT"))
    if err==nil{
        app.Server.HTTP.Port=port
    }

    app.Start()
}