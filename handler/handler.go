package handler

import (
	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/aryanA101a/legoshichat-backend/store"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
)

type handler struct {
	auth store.AuthStore
}

func New(a store.AuthStore) handler {
	return handler{auth: a}
}

func (h handler) CreateAccount(ctx *gofr.Context) (interface{}, error) {
	var accountRequest model.CreateAccountRequest

	if err := ctx.Bind(&accountRequest); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := h.auth.CreateAccount(ctx, accountRequest)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

