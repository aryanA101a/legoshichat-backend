package store

import (
	"database/sql"

	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/aryanA101a/legoshichat-backend/model"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
}

type AuthStore interface {
	CreateAccount(ctx *gofr.Context, account model.Account) (*model.User, error)
	Login(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.User, error)
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

func (auth) Login(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.User, error) {
	var user model.User
	var passwordHash string

	err := ctx.DB().QueryRowContext(ctx, " SELECT id,name,phoneNumber,password FROM accounts where phoneNumber=$1", loginRequest.PhoneNumber).
		Scan(&user.ID, &user.Name, &user.PhoneNumber, &passwordHash)
	if err != nil {
		ctx.Logger.Info("err: ", err.Error())
		if err == sql.ErrNoRows {
			return nil,e.HttpStatusError(401,"User does not exists") 
		}
		return nil, e.HttpStatusError(500,"")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginRequest.Password))
	if err != nil {
		return nil, e.HttpStatusError(401,"Invalid password")
	}

	return &user, nil

}

func (auth) createAccountsTable(db *datastore.SQLClient) error {
	query := "CREATE TABLE IF NOT EXISTS accounts" +
		" (id uuid PRIMARY KEY , name varchar(100) NOT NULL, phoneNumber bigint UNIQUE NOT NULL, password varchar(100) NOT NULL);"
	_, err := db.Exec(query)
	return err
}
