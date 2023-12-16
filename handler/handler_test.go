package handler

import (
	"bytes"
	"context"
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

func (mockAuthStore) AccountExists(ctx *gofr.Context, userId string) (bool, error) {

	return false, nil
}

func (mockAuthStore) GetUserIdByPhoneNumber(ctx *gofr.Context, phoneNumber uint64) (*string, error) {

	return nil, nil
}

func (createAccountTCmockAuthCreator) NewAccount(name string, phoneNumber uint64, password string) (*model.Account, error) {
	fmt.Println("t3mockAuthCreator called")
	return nil, e.NewError("")
}

func (createAccountTCmockAuthCreator) NewJWT(ctx *gofr.Context, userId string) (string, error) {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOiIyMDIzLTEyLTE0VDE1OjA5OjA1LjY2MzE5NjA1NiswNTozMCIsImlkIjoiY2E0MGZkMmItZjlhMi00NWQ5LWE2ZjctNDFjZGFiNWVjMDFlIn0.2Wp_KL-rM7RMQmWCCLvtHrG5M-zVuGFO-7PBLzDAbJ4", nil
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
		err:  e.NewError(""),
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

	runTest(t, successfulTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: NewCreator()})
	runTest(t, invalidRequestBodyTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: NewCreator()})
	runTest(t, createAccountErrorTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: createAccountTCmockAuthCreator{}})
	runTest(t, jwtErrorTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: jwtErrorTCmockAuthCreator{}})

}

func runTest(t *testing.T, tc testCase, app *gofr.Gofr, h Handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)

	result, err := h.HandleCreateAccount(ctx)

	fmt.Println("errr", tc.expected, err)

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

func (fetchAccountErrorTcmockAuthStore) AccountExists(ctx *gofr.Context, userId string) (bool, error) {

	return false, nil
}

func (fetchAccountErrorTcmockAuthStore) GetUserIdByPhoneNumber(ctx *gofr.Context, phoneNumber uint64) (*string, error) {

	return nil, nil
}

type userNotExistsTCmockAuthStore struct{}

func (userNotExistsTCmockAuthStore) CreateAccount(ctx *gofr.Context, account model.Account) (*model.User, error) {

	return &model.User{ID: account.ID, Name: account.Name, PhoneNumber: account.PhoneNumber}, nil

}

func (userNotExistsTCmockAuthStore) FetchAccount(ctx *gofr.Context, loginRequest model.LoginRequest) (*model.Account, error) {

	return nil, sql.ErrNoRows
}

func (userNotExistsTCmockAuthStore) AccountExists(ctx *gofr.Context, userId string) (bool, error) {

	return false, nil
}

func (userNotExistsTCmockAuthStore) GetUserIdByPhoneNumber(ctx *gofr.Context, phoneNumber uint64) (*string, error) {

	return nil, nil
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


	runLoginTest(t, successfulLoginTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: NewCreator()})
	runLoginTest(t, invalidRequestBodyTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: NewCreator()})
	runLoginTest(t, fetchAccountErrorTC, app, Handler{Auth: fetchAccountErrorTcmockAuthStore{}, AuthCreator: NewCreator()})
	runLoginTest(t, userNotExistsTC, app, Handler{Auth: userNotExistsTCmockAuthStore{}, AuthCreator: NewCreator()})
	runLoginTest(t, invalidPasswordTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: NewCreator()})
	runLoginTest(t, jwtErrorTC, app, Handler{Auth: newMockAuthStore(), AuthCreator: jwtErrorTCmockAuthCreator{}})
}

func runLoginTest(t *testing.T, tc testCase, app *gofr.Gofr, h Handler) {
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

type successfulTCMessageStore struct{}

func (successfulTCMessageStore) AddMessage(ctx *gofr.Context, message model.Message) error {
	return nil
}

func (successfulTCMessageStore) GetMessage(ctx *gofr.Context, userId, messageId string) (*model.Message, error) {
	return nil, nil
}

func (successfulTCMessageStore) UpdateMessage(ctx *gofr.Context, userId, messageId, updatedContent string) (*model.Message, error) {
	return nil, nil
}

func (successfulTCMessageStore) DeleteMessage(ctx *gofr.Context, userId, messageId string) error {
	return nil
}

func (successfulTCMessageStore) GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error) {
	return nil, nil
}

type addMessageErrorTCMessageStore struct{}

func (addMessageErrorTCMessageStore) AddMessage(ctx *gofr.Context, message model.Message) error {
	return e.NewError("")
}

func (addMessageErrorTCMessageStore) GetMessage(ctx *gofr.Context, userId, messageId string) (*model.Message, error) {
	return nil, nil
}

func (addMessageErrorTCMessageStore) UpdateMessage(ctx *gofr.Context, userId, messageId, updatedContent string) (*model.Message, error) {
	return nil, nil
}

func (addMessageErrorTCMessageStore) DeleteMessage(ctx *gofr.Context, userId, messageId string) error {
	return nil
}

func (addMessageErrorTCMessageStore) GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error) {
	return nil, nil
}

func TestHandleSendMessageByID(t *testing.T) {
	app := gofr.New()

	validInputTC := testCaseSendMessage{
		desc:     "send message by ID success",
		body:     []byte(`{"recipientID":"123","content":"Hello, World!"}`),
		userID:   "testUserID",
		expected: types.Raw{},
		err:      nil,
	}

	invalidInputTC := testCaseSendMessage{
		desc: "invalid input",
		body: []byte(`invalidjson`),
		err:  e.HttpStatusError(400, "Invalid inputs or missing required fields"),
	}

	messageStoreErrorTC := testCaseSendMessage{
		desc:   "message store error",
		body:   []byte(`{"recipientID":"123","content":"Hello, World!"}`),
		userID: "testUserID",
		err:    e.HttpStatusError(500, "message store error"),
	}


	runSendMessageTest(t, validInputTC, app, Handler{Message: successfulTCMessageStore{}})
	runSendMessageTest(t, invalidInputTC, app, Handler{Message: successfulTCMessageStore{}})
	runSendMessageTest(t, messageStoreErrorTC, app, Handler{Message: addMessageErrorTCMessageStore{}})
}

type testCaseSendMessage struct {
	desc     string
	body     []byte
	userID   string
	expected interface{}
	err      error
}

func runSendMessageTest(t *testing.T, tc testCaseSendMessage, app *gofr.Gofr, h Handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)
	*&ctx.Context = context.WithValue(context.Background(), "userId", "someUserId")

	result, err := h.HandleSendMessageByID(ctx)

	if tc.err != nil {
		assert.Error(t, err, "TEST: %s: unexpected error", tc.desc)
	} else {
		assert.NoError(t, err, "TEST: Unexpected Error: %s", tc.desc)
		assert.IsType(t, tc.expected, result, "TEST: %s: unexpected result type", tc.desc)
	}
}
