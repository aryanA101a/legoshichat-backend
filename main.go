package main

import (
	"strconv"

	"github.com/aryanA101a/legoshichat-backend/handler"
	"github.com/aryanA101a/legoshichat-backend/store"
	"gofr.dev/pkg/gofr"
)

func main() {
    app := gofr.New()

    s := store.NewAuthStore(app.DB())
	h := handler.New(s)
	
    
    app.POST("/create-account", h.CreateAccount)
    
    port,err:=strconv.Atoi(app.Config.Get("HTTP_PORT"))
    if err==nil{
        app.Server.HTTP.Port=port
    }

    app.Start()
}