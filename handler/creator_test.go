package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
)

func TestNewAccount(t *testing.T) {

	creator := NewCreator()

	account, err := creator.NewAccount("TestUser", 1234567890, "securepassword")

	assert.NoError(t, err, "Unexpected error during account creation")
	assert.NotNil(t, account, "Expected a non-nil account object")
	assert.NotEmpty(t, account.ID, "Expected a non-empty account ID")
	assert.Equal(t, "TestUser", account.Name, "Mismatch in account name")
	assert.Equal(t, uint64(1234567890), account.PhoneNumber, "Mismatch in account phone number")
}

func TestNewJWT(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())

	creator := NewCreator()

	userID := "test-user-id"

	token, err := creator.NewJWT(ctx, userID)

	assert.NoError(t, err, "Unexpected error during JWT creation")
	assert.NotEmpty(t, token, "Expected a non-empty JWT token")
}