package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/request"
	"gofr.dev/pkg/gofr/responder"
	"gofr.dev/pkg/gofr/types"

	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/aryanA101a/legoshichat-backend/store"
	"github.com/stretchr/testify/assert"
)

type mockAuthStore struct {
}

type createAccountTCmockAuthCreator struct {
}
type jwtErrorTCmockAuthCreator struct {
}

func newMockAuthStore() store.AuthStore {
	return mockAuthStore{}
}


func (mockAuthStore) CreateAccount(ctx *gofr.Context, account model.Account) (*model.User, error) {

	return &model.User{ID: account.ID, Name: account.Name, PhoneNumber: account.PhoneNumber}, nil

}

func (mockAuthStore) FetchAccount(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.Account, error) {

	return &model.Account{Password: "$2a$10$hlaojbP2Ix5ZLFa64fPmgO3tyvH.jdBDVtq8veZyQMOWohcAndlpS"}, nil
}

func (createAccountTCmockAuthCreator) NewAccount(name string, phoneNumber uint64, password string) (*model.Account, error) {
fmt.Println("t3mockAuthCreator called")
	return nil, e.NewError("")
}

func (createAccountTCmockAuthCreator) NewJWT(ctx *gofr.Context, userId string) (string, error) {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOiIyMDIzLTEyLTE0VDE1OjA5OjA1LjY2MzE5NjA1NiswNTozMCIsImlkIjoiY2E0MGZkMmItZjlhMi00NWQ5LWE2ZjctNDFjZGFiNWVjMDFlIn0.2Wp_KL-rM7RMQmWCCLvtHrG5M-zVuGFO-7PBLzDAbJ4",nil
}

func (jwtErrorTCmockAuthCreator) NewAccount(name string, phoneNumber uint64, password string) (*model.Account, error) {

	return &model.Account{}, nil
}

func (jwtErrorTCmockAuthCreator) NewJWT(ctx *gofr.Context, userId string) (string, error) {
	return "", e.NewError("")
}

type testCase struct {
	desc     string
	body     []byte
	expected interface{}
	err      error
}

func TestHandleCreateAccount(t *testing.T) {

	app := gofr.New()

	successfulTC := testCase{
		desc:     "create account success",
		body:     []byte(`{"name":"TestUser4","phoneNumber":4234324890,"password":"securepassword"}`),
		expected: types.Raw{Data: model.AuthResponse{}},
		err:      nil,
	}

	invalidRequestBodyTC := testCase{
		desc: "invalid request body",
		body: []byte(`invalidjson`),
		err: e.NewError(""),
	}

	createAccountErrorTC := testCase{
		desc: "create account error",
		body: []byte(`{"name":"TestUser4","phoneNumber":4234324890,"password":"securepassword"}`),
		err:  e.NewError(""),
	}

	jwtErrorTC := testCase{
		desc: "create account JWT error",
		body: []byte(`{"name":"TestUser","phoneNumber":"1234567890","password":"securepassword"}`),
		err:  e.NewError(""),
	}

	runTest(t, successfulTC, app, New(newMockAuthStore(), NewCreator()))
	runTest(t, invalidRequestBodyTC, app, New(newMockAuthStore(), NewCreator()))
	runTest(t, createAccountErrorTC, app, New(newMockAuthStore(), createAccountTCmockAuthCreator{}))
	runTest(t, jwtErrorTC, app, New(newMockAuthStore(), jwtErrorTCmockAuthCreator{}))


}

func runTest(t *testing.T, tc testCase, app *gofr.Gofr, h handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)

	result, err := h.HandleCreateAccount(ctx)

	fmt.Println("errr", tc.expected,err)

	if tc.err != nil {
		assert.Error(t, err, "TEST: %s: unexpected error", tc.desc)
	} else {
		assert.NoError(t, err, "TEST: Unexpected Error: %s", tc.desc)
	}

	assert.IsType(t, tc.expected, result, "TEST: %s: unexpected result type", tc.desc)
}


type fetchAccountErrorTcmockAuthStore struct{}

func (fetchAccountErrorTcmockAuthStore) CreateAccount(ctx *gofr.Context, account model.Account) (*model.User, error) {

	return &model.User{ID: account.ID, Name: account.Name, PhoneNumber: account.PhoneNumber}, nil

}

func (fetchAccountErrorTcmockAuthStore) FetchAccount(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.Account, error) {

	return nil, e.NewError("")
}

type userNotExistsTCmockAuthStore struct{}

func (userNotExistsTCmockAuthStore) CreateAccount(ctx *gofr.Context, account model.Account) (*model.User, error) {

	return &model.User{ID: account.ID, Name: account.Name, PhoneNumber: account.PhoneNumber}, nil

}

func (userNotExistsTCmockAuthStore) FetchAccount(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.Account, error) {

	return nil, sql.ErrNoRows
}
func TestHandleLogin(t *testing.T) {
	app := gofr.New()

	successfulLoginTC := testCase{
		desc: "login success",
		body: []byte(`{"phoneNumber":1234567890,"password":"unittestpass"}`),
		expected: types.Raw{
			Data: model.AuthResponse{
				// User:  model.User{ID: "testUserID", Name: "TestUser", PhoneNumber: 4234324890},
				// Token: "testToken",
			},
		},
		err: nil,
	}

	invalidRequestBodyTC := testCase{
		desc: "invalid request body",
		body: []byte(`invalidjson`),
		err:  e.NewError(""),
	}

	fetchAccountErrorTC := testCase{
		desc: "fetch account error",
		body: []byte(`{"phoneNumber":1234567890,"password":"securepassword"}`),
		err:  e.HttpStatusError(500, ""),
	}

	userNotExistsTC := testCase{
		desc: "user does not exist",
		body: []byte(`{"phoneNumber":1234567892,"password":"securepassword"}`),
		err:  e.HttpStatusError(401, "User does not exist"),
	}

	invalidPasswordTC := testCase{
		desc: "invalid password",
		body: []byte(`{"phoneNumber":1234567890,"password":"wrongpassword"}`),
		err:  e.HttpStatusError(401, "Invalid password"),
	}

	jwtErrorTC := testCase{
		desc: "login JWT error",
		body: []byte(`{"phoneNumber":1234567890,"password":"unittestpass"}`),
		err:  e.NewError(""),
	}

	// Add more test cases as needed

	runLoginTest(t, successfulLoginTC, app, New(newMockAuthStore(), NewCreator(), ))
	runLoginTest(t, invalidRequestBodyTC, app, New(newMockAuthStore(), NewCreator()))
	runLoginTest(t, fetchAccountErrorTC, app, New(fetchAccountErrorTcmockAuthStore{}, NewCreator(), ))
	runLoginTest(t, userNotExistsTC, app, New(userNotExistsTCmockAuthStore{}, NewCreator(),))
	runLoginTest(t, invalidPasswordTC, app, New(newMockAuthStore(), NewCreator(),))
	runLoginTest(t, jwtErrorTC, app, New(newMockAuthStore(), jwtErrorTCmockAuthCreator{},))
}


func runLoginTest(t *testing.T, tc testCase, app *gofr.Gofr, h handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)

	result, err := h.HandleLogin(ctx)

	fmt.Println("errr", tc.expected, err)

	if tc.err != nil {
		assert.Error(t, err, "TEST: %s: unexpected error", tc.desc)
	} else {
		assert.NoError(t, err, "TEST: Unexpected Error: %s", tc.desc)
	}

	assert.IsType(t, tc.expected, result, "TEST: %s: unexpected result type", tc.desc)
}