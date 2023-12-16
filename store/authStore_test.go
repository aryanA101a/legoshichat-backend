package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
)


func TestCreateAccount(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	defer db.Close()

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(sqlmock.AnyArg(), "TestUser", uint64(1234567890), "securepassword").
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	authStore := NewAuthStore(ctx.DB())

	successfulAccount := model.Account{
		ID:          "test-account-id",
		Name:        "TestUser",
		PhoneNumber: 1234567890,
		Password:    "securepassword",
	}

	createdUser, err := authStore.CreateAccount(ctx, successfulAccount)

	assert.NoError(t, err, "Unexpected error during successful account creation")
	assert.NotNil(t, createdUser, "Expected a non-nil user object after successful account creation")
	assert.Equal(t, successfulAccount.ID, createdUser.ID, "Mismatch in created user ID")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	failedAccount := model.Account{
		ID:          "test-account-id-fail",
		Name:        "TestUserFail",
		PhoneNumber: 9876543210,
		Password:    "passwordfail",
	}

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(sqlmock.AnyArg(), "TestUserFail", uint64(9876543210), "passwordfail").
		WillReturnError(fmt.Errorf(""))

	createdUserFail, errFail := authStore.CreateAccount(ctx, failedAccount)

	assert.Error(t, errFail, "Expected an error during failed account creation")
	assert.Nil(t, createdUserFail, "Expected a nil user object after failed account creation")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}


func TestFetchAccount(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	defer db.Close()

	mock.ExpectQuery("SELECT id,name,phoneNumber,password FROM accounts").
		WithArgs(uint64(1234567890)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phoneNumber", "password"}).
			AddRow("test-account-id", "TestUser", uint64(1234567890), "securepassword"))

	authStore := NewAuthStore(ctx.DB())

	loginRequest := model.LoginRequest{
		PhoneNumber: 1234567890,
	}

	fetchedAccount, err := authStore.FetchAccount(ctx, loginRequest)

	assert.NoError(t, err, "Unexpected error during successful account fetch")
	assert.NotNil(t, fetchedAccount, "Expected a non-nil account object after successful fetch")
	assert.Equal(t, "test-account-id", fetchedAccount.ID, "Mismatch in fetched account ID")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id,name,phoneNumber,password FROM accounts").
		WithArgs(uint64(9876543210)).
		WillReturnError(fmt.Errorf(""))

	failedLoginRequest := model.LoginRequest{
		PhoneNumber: 9876543210,
	}

	fetchedAccountFail, errFail := authStore.FetchAccount(ctx, failedLoginRequest)

	assert.Error(t, errFail, "Expected an error during failed account fetch")
	assert.Nil(t, fetchedAccountFail, "Expected a nil account object after failed fetch")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAccountExists(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	defer db.Close()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("test-user-id").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	authStore := NewAuthStore(ctx.DB())

	exists, err := authStore.AccountExists(ctx, "test-user-id")

	assert.NoError(t, err, "Unexpected error during account existence check")
	assert.True(t, exists, "Expected the account to exist")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("non-existent-user-id").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	notExists, err := authStore.AccountExists(ctx, "non-existent-user-id")

	assert.NoError(t, err, "Unexpected error during account existence check")
	assert.False(t, notExists, "Expected the account to not exist")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("error-user-id").
		WillReturnError(fmt.Errorf(""))

	_, errFail := authStore.AccountExists(ctx, "error-user-id")

	assert.Error(t, errFail, "Expected an error during failed account existence check")
}

func TestGetUserIdByPhoneNumber(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	defer db.Close()

	mock.ExpectQuery("SELECT id FROM accounts WHERE phoneNumber").
		WithArgs(uint64(1234567890)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("test-user-id"))

	authStore := NewAuthStore(ctx.DB())

	userID, err := authStore.GetUserIdByPhoneNumber(ctx, uint64(1234567890))

	assert.NoError(t, err, "Unexpected error during phone number to user ID retrieval")
	assert.NotNil(t, userID, "Expected a non-nil user ID after successful retrieval")
	assert.Equal(t, "test-user-id", *userID, "Mismatch in retrieved user ID")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id FROM accounts WHERE phoneNumber").
		WithArgs(uint64(9876543210)).
		WillReturnError(fmt.Errorf(""))

	_, errFail := authStore.GetUserIdByPhoneNumber(ctx, uint64(9876543210))

	assert.Error(t, errFail, "Expected an error during failed phone number to user ID retrieval")
}

func TestCreateAccountsTable(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS accounts").
		WillReturnResult(sqlmock.NewResult(0, 0)).
		WillReturnError(nil)


	err = createAccountsTable(ctx.DB())

	assert.NoError(t, err, "Unexpected error during successful table creation")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS accounts").
		WillReturnError(fmt.Errorf(""))

	errFail := createAccountsTable(ctx.DB())

	assert.Error(t, errFail, "Expected an error during failed table creation")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}