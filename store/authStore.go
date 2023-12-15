package store

import (
	"github.com/aryanA101a/legoshichat-backend/model"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
)

type auth struct {
}

type AuthStore interface {
	CreateAccount(ctx *gofr.Context, account model.Account) (*model.User, error)
	FetchAccount(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.Account, error)
	AccountExists(ctx *gofr.Context, userId string) (bool, error)
	GetUserIdByPhoneNumber(ctx *gofr.Context, phoneNumber uint64) (*string, error)
}

func NewAuthStore(db *datastore.SQLClient) AuthStore {
	a := auth{}
	a.init(db)
	return a
}

func (a auth) init(db *datastore.SQLClient) {
	a.createAccountsTable(db)
}

func (auth) CreateAccount(ctx *gofr.Context, account model.Account) (*model.User, error) {

	_, err := ctx.DB().ExecContext(ctx, "INSERT INTO accounts (id,name,phoneNumber,password) VALUES ($1,$2,$3,$4)", account.ID, account.Name, account.PhoneNumber, account.Password)
	if err != nil {
		return nil, err
	}

	return &model.User{ID: account.ID, Name: account.Name, PhoneNumber: account.PhoneNumber}, nil
}

func (auth) FetchAccount(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.Account, error) {
	var account model.Account

	err := ctx.DB().QueryRowContext(ctx, " SELECT id,name,phoneNumber,password FROM accounts where phoneNumber=$1", loginRequest.PhoneNumber).
		Scan(&account.ID, &account.Name, &account.PhoneNumber, &account.Password)
	if err != nil {
		return nil, err
	}

	return &account, nil

}

func (a auth) AccountExists(ctx *gofr.Context, userId string) (bool, error) {
	var exists bool
	err := ctx.DB().QueryRowContext(ctx, "SELECT EXISTS(SELECT id FROM accounts WHERE id=$1)", userId).
		Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a auth) GetUserIdByPhoneNumber(ctx *gofr.Context, phoneNumber uint64) (*string, error) {
	var id string
	err := ctx.DB().QueryRowContext(ctx, "SELECT id FROM accounts WHERE phoneNumber=$1", phoneNumber).
		Scan(&id)

	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (auth) createAccountsTable(db *datastore.SQLClient) error {
	query := "CREATE TABLE IF NOT EXISTS accounts" +
		" (id uuid PRIMARY KEY , name varchar(100) NOT NULL, phoneNumber bigint UNIQUE NOT NULL, password varchar(100) NOT NULL);"
	_, err := db.Exec(query)
	return err
}
