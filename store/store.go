package store

import (
	"github.com/aryanA101a/legoshichat-backend/model"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
)

type auth struct {
}

type AuthStore interface {
	CreateAccount(ctx *gofr.Context, accountRequest model.CreateAccountRequest) (*model.User, error)
}

func NewAuthStore(db *datastore.SQLClient) AuthStore {
	a := auth{}
	a.init(db)
	return a
}

func (a auth) init(db *datastore.SQLClient) {
	a.createAccountsTable(db)
}

func (auth) CreateAccount(ctx *gofr.Context, accountRequest model.CreateAccountRequest) (*model.User, error) {

	newAccount, err := model.NewAccount(accountRequest.Name, accountRequest.PhoneNumber, accountRequest.Password)
	if err != nil {
		return nil, err
	}
	
	_, err = ctx.DB().ExecContext(ctx, "INSERT INTO accounts (id,name,phoneNumber,password) VALUES ($1,$2,$3,$4)", newAccount.ID, newAccount.Name, newAccount.PhoneNumber, newAccount.Password)
	if err != nil {
		return nil, err
	}

	return &model.User{ID: newAccount.ID, Name: newAccount.Name, PhoneNumber: newAccount.PhoneNumber}, nil
}

func (auth) createAccountsTable(db *datastore.SQLClient) error {
	query := "CREATE TABLE IF NOT EXISTS accounts" +
		" (id uuid PRIMARY KEY , name varchar(100) NOT NULL, phoneNumber bigint UNIQUE NOT NULL, password varchar(100) NOT NULL);"
	_, err := db.Exec(query)
	return err
}
