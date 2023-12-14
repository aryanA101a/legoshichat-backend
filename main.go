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
	handler := handler.New(authStore)
	
    
    app.POST("/create-account", handler.HandleCreateAccount)
    app.POST("/login", handler.HandleLogin)
    
    port,err:=strconv.Atoi(app.Config.Get("HTTP_PORT"))
    if err==nil{
        app.Server.HTTP.Port=port
    }

    app.Start()
}