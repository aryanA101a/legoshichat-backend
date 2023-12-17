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
var r string=""
	return &r, nil
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
	return &model.Message{}, nil
}

func (successfulTCMessageStore) DeleteMessage(ctx *gofr.Context, userId, messageId string) error {
	return nil
}

func (successfulTCMessageStore) GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error) {
	return nil, nil
}

type errorTCMessageStore struct{}

func (errorTCMessageStore) AddMessage(ctx *gofr.Context, message model.Message) error {
	return e.NewError("")
}

func (errorTCMessageStore) GetMessage(ctx *gofr.Context, userId, messageId string) (*model.Message, error) {
	return nil, sql.ErrNoRows
}

func (errorTCMessageStore) UpdateMessage(ctx *gofr.Context, userId, messageId, updatedContent string) (*model.Message, error) {
	return nil, sql.ErrNoRows
}

func (errorTCMessageStore) DeleteMessage(ctx *gofr.Context, userId, messageId string) error {
	return sql.ErrNoRows
}

func (errorTCMessageStore) GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error) {
	return nil, nil
}

type messageStoreErrorTCMessageStore struct{}

func (messageStoreErrorTCMessageStore) AddMessage(ctx *gofr.Context, message model.Message) error {
	return e.NewError("")
}

func (messageStoreErrorTCMessageStore) GetMessage(ctx *gofr.Context, userId, messageId string) (*model.Message, error) {
	return nil, e.NewError("")
}

func (messageStoreErrorTCMessageStore) UpdateMessage(ctx *gofr.Context, userId, messageId, updatedContent string) (*model.Message, error) {
	return nil, e.NewError("")
}

func (messageStoreErrorTCMessageStore) DeleteMessage(ctx *gofr.Context, userId, messageId string) error {
	return e.NewError("")
}

func (messageStoreErrorTCMessageStore) GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error) {
	return nil, nil
}

type authorizationErrorTCMessageStore struct{}

func (authorizationErrorTCMessageStore) AddMessage(ctx *gofr.Context, message model.Message) error {
	return e.NewError("")
}

func (authorizationErrorTCMessageStore) GetMessage(ctx *gofr.Context, userId, messageId string) (*model.Message, error) {
	return nil, e.NewError("You are not authorized to see that message")
}

func (authorizationErrorTCMessageStore) UpdateMessage(ctx *gofr.Context, userId, messageId, updatedContent string) (*model.Message, error) {
	return nil, e.NewError("You are not authorized to see that message")
}

func (authorizationErrorTCMessageStore) DeleteMessage(ctx *gofr.Context, userId, messageId string) error {
	return e.NewError("You are not authorized to see that message")
}

func (authorizationErrorTCMessageStore) GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error) {
	return nil, nil
}

type mockFriendStore struct{}

func (f mockFriendStore) AddFriend(ctx *gofr.Context, senderId, recieverId string) error {
	return nil
}

func (f mockFriendStore) GetFriends(ctx *gofr.Context, userId string) (*[]model.User, error) {
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

	runSendMessageTest(t, validInputTC, app, Handler{Message: successfulTCMessageStore{}, Friend: mockFriendStore{}})
	runSendMessageTest(t, invalidInputTC, app, Handler{Message: successfulTCMessageStore{}, Friend: mockFriendStore{}})
	runSendMessageTest(t, messageStoreErrorTC, app, Handler{Message: errorTCMessageStore{}, Friend: mockFriendStore{}})
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

func TestHandleSendMessageByPhoneNumber(t *testing.T) {
	app := gofr.New()

	validInputTC := testCaseSendMessageByPhoneNumber{
		desc:     "send message by phone number success",
		body:     []byte(`{"recipientPhoneNumber":1234567890,"content":"Hello, World!"}`),
		userID:   "testUserID",
		expected: types.Raw{},
		err:      nil,
	}

	invalidInputTC := testCaseSendMessageByPhoneNumber{
		desc: "invalid input",
		body: []byte(`invalidjson`),
		err:  e.HttpStatusError(400, "Invalid inputs or missing required fields -json: unknown field \"invalidjson\""),
	}

	messageStoreErrorTC := testCaseSendMessageByPhoneNumber{
		desc:   "message store error",
		body:   []byte(`{"recipientPhoneNumber":1234567890,"content":"Hello, World!"}`),
		userID: "testUserID",
		err:    e.HttpStatusError(500, "message store error"),
	}

	recipientNotFoundTC := testCaseSendMessageByPhoneNumber{
		desc:   "recipient not found",
		body:   []byte(`{"recipientPhoneNumber":1234567890,"content":"Hello, World!"}`),
		userID: "testUserID",
		err:    e.HttpStatusError(404, "Recipient does not exists"),
	}

	runSendMessageByPhoneNumberTest(t, validInputTC, app, Handler{Message: successfulTCMessageStore{}, Friend: mockFriendStore{}, Auth: newMockAuthStore()})
	runSendMessageByPhoneNumberTest(t, invalidInputTC, app, Handler{Message: successfulTCMessageStore{}, Friend: mockFriendStore{}, Auth:newMockAuthStore()})
	runSendMessageByPhoneNumberTest(t, messageStoreErrorTC, app, Handler{Message: errorTCMessageStore{}, Friend: mockFriendStore{}, Auth:newMockAuthStore()})
	runSendMessageByPhoneNumberTest(t, recipientNotFoundTC, app, Handler{Message: successfulTCMessageStore{}, Friend: mockFriendStore{}, Auth: userNotExistsTCmockAuthStore{}})
}

type testCaseSendMessageByPhoneNumber struct {
	desc     string
	body     []byte
	userID   string
	expected interface{}
	err      error
}

func runSendMessageByPhoneNumberTest(t *testing.T, tc testCaseSendMessageByPhoneNumber, app *gofr.Gofr, h Handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)
	*&ctx.Context = context.WithValue(context.Background(), "userId", "someUserId")

	result, err := h.HandleSendMessageByPhoneNumber(ctx)

	if tc.err != nil {
		assert.Error(t, err, "TEST: %s: unexpected error", tc.desc)
	} else {
		assert.NoError(t, err, "TEST: Unexpected Error: %s", tc.desc)
		assert.IsType(t, tc.expected, result, "TEST: %s: unexpected result type", tc.desc)
	}
}

func TestHandleGetMessage(t *testing.T) {
	app := gofr.New()

	validInputTC := testCaseGetDeleteMessage{
		desc:     "get message success",
		messageID: "testMessageID",
		userID:    "testUserID",
		expected:  types.Raw{},
		err:       nil,
	}

	missingParameterTC := testCaseGetDeleteMessage{
		desc:     "missing parameter",
		messageID: "",
		userID:    "testUserID",
		err:       e.HttpStatusError(400, "Missing Parameter messageId"),
	}

	messageNotFoundErrorTC := testCaseGetDeleteMessage{
		desc:     "message not found",
		messageID: "nonexistentMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(404, "Message does not exists"),
	}
	messageStoreErrorTC := testCaseGetDeleteMessage{
		desc:      "message store error",
		messageID: "testMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(500, "message store error"),
	}
	authorizationErrorTC := testCaseGetDeleteMessage{
		desc:     "authorization error",
		messageID: "unauthorizedMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(403, "You are not authorized to see that message"),
	}

	runGetMessageTest(t, validInputTC, app, Handler{Message: successfulTCMessageStore{}})
	runGetMessageTest(t, missingParameterTC, app, Handler{Message: successfulTCMessageStore{}})
	runGetMessageTest(t, messageNotFoundErrorTC, app, Handler{Message: errorTCMessageStore{}})
	runGetMessageTest(t, messageStoreErrorTC, app, Handler{Message: messageStoreErrorTCMessageStore{}})
	runGetMessageTest(t, authorizationErrorTC, app, Handler{Message: authorizationErrorTCMessageStore{}})
}

type testCaseGetDeleteMessage struct {
	desc      string
	messageID string
	userID    string
	expected  interface{}
	err       error
}

func runGetMessageTest(t *testing.T, tc testCaseGetDeleteMessage, app *gofr.Gofr, h Handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy",nil)

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)
	*&ctx.Context = context.WithValue(context.Background(), "userId", "someUserId")
	
	if tc.desc!="missing parameter"{
		ctx.SetPathParams(map[string]string{"id":"someOtherUserId"})
	}

	result, err := h.HandleGetMessage(ctx)

	if tc.err != nil {
		assert.Error(t, err, "TEST: %s: unexpected error", tc.desc)
	} else {
		assert.NoError(t, err, "TEST: Unexpected Error: %s", tc.desc)
		assert.IsType(t, tc.expected, result, "TEST: %s: unexpected result type", tc.desc)
	}
}


type testCasePutMessage struct {
	desc      string
	body      []byte
	messageID string
	userID    string
	expected  interface{}
	err       error
}
func TestHandlePutMessage(t *testing.T) {
	app := gofr.New()


	validInputTC := testCasePutMessage{
		desc:     "update message success",
		body:     []byte(`{"content":"Updated content"}`),
		messageID: "testMessageID",
		userID:    "testUserID",
		expected:  types.Raw{},
		err:       nil,
	}
	invalidBodyTC := testCasePutMessage{
		desc:      "invalid request body",
		body:      []byte(`invalidjson`),
		messageID: "testMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(400, "Invalid inputs or missing required fields -json: cannot unmarshal string into Go struct field UpdateMessageRequest.content of type model.Content"),
	}

	missingParameterTC := testCasePutMessage{
		desc:      "missing parameter",
		body:      []byte(`{"content":"Updated content"}`),
		messageID: "",
		userID:    "testUserID",
		err:       e.HttpStatusError(400, "Missing Parameter messageId"),
	}

	messageNotFoundErrorTC := testCasePutMessage{
		desc:      "message not found",
		body:      []byte(`{"content":"Updated content"}`),
		messageID: "nonexistentMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(404, "Message does not exists"),
	}
	messageStoreErrorTC := testCasePutMessage{
		desc:      "message store error",
		body:      []byte(`{"content":"Updated content"}`),
		messageID: "testMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(500, "message store error"),
	}
	authorizationErrorTC := testCasePutMessage{
		desc:      "authorization error",
		body:      []byte(`{"content":"Updated content"}`),
		messageID: "unauthorizedMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(403, "You are not authorized to update that message"),
	}

	runPutMessageTest(t, validInputTC, app, Handler{Message: successfulTCMessageStore{}})
	runPutMessageTest(t, invalidBodyTC, app, Handler{Message: successfulTCMessageStore{}})
	runPutMessageTest(t, missingParameterTC, app, Handler{Message: successfulTCMessageStore{}})
	runPutMessageTest(t, messageNotFoundErrorTC, app, Handler{Message: errorTCMessageStore{}})
	runPutMessageTest(t, messageStoreErrorTC, app, Handler{Message: messageStoreErrorTCMessageStore{}})
	runPutMessageTest(t, authorizationErrorTC, app, Handler{Message: authorizationErrorTCMessageStore{}})

}

func runPutMessageTest(t *testing.T, tc testCasePutMessage, app *gofr.Gofr, h Handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)
	*&ctx.Context = context.WithValue(context.Background(), "userId", "someUserId")
	
	if tc.desc!="missing parameter"{
		ctx.SetPathParams(map[string]string{"id":"someOtherUserId"})
	}

	result, err := h.HandlePutMessage(ctx)

	if tc.err != nil {
		assert.Error(t, err, "TEST: %s: unexpected error", tc.desc)
	} else {
		assert.NoError(t, err, "TEST: Unexpected Error: %s", tc.desc)
		assert.IsType(t, tc.expected, result, "TEST: %s: unexpected result type", tc.desc)
	}
}


func TestHandleDeleteMessage(t *testing.T) {
	app := gofr.New()

	validInputTC := testCaseGetDeleteMessage{
		desc:     "delete message success",
		messageID: "testMessageID",
		userID:    "testUserID",
		err:       nil,
	}

	missingParameterTC := testCaseGetDeleteMessage{
		desc:     "missing parameter",
		messageID: "",
		userID:    "testUserID",
		err:       e.HttpStatusError(400, "Missing Parameter messageId"),
	}

	messageNotFoundErrorTC := testCaseGetDeleteMessage{
		desc:     "message not found",
		messageID: "nonexistentMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(404, "Message does not exists"),
	}
	messageStoreErrorTC := testCaseGetDeleteMessage{
		desc:      "message store error",
		messageID: "testMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(500, "message store error"),
	}
	authorizationErrorTC := testCaseGetDeleteMessage{
		desc:     "authorization error",
		messageID: "unauthorizedMessageID",
		userID:    "testUserID",
		err:       e.HttpStatusError(403, "You are not authorized to see that message"),
	}

	runDeleteMessageTest(t, validInputTC, app, Handler{Message: successfulTCMessageStore{}})
	runDeleteMessageTest(t, missingParameterTC, app, Handler{Message: successfulTCMessageStore{}})
	runDeleteMessageTest(t, messageNotFoundErrorTC, app, Handler{Message: errorTCMessageStore{}})
	runDeleteMessageTest(t, messageStoreErrorTC, app, Handler{Message: messageStoreErrorTCMessageStore{}})
	runDeleteMessageTest(t, authorizationErrorTC, app, Handler{Message: authorizationErrorTCMessageStore{}})
}



func runDeleteMessageTest(t *testing.T, tc testCaseGetDeleteMessage, app *gofr.Gofr, h Handler) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://dummy",nil)

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)
	*&ctx.Context = context.WithValue(context.Background(), "userId", "someUserId")
	
	if tc.desc!="missing parameter"{
		ctx.SetPathParams(map[string]string{"id":"someOtherUserId"})
	}

	result, err := h.HandleDeleteMessage(ctx)

	if tc.err != nil {
		assert.Error(t, err, "TEST: %s: unexpected error", tc.desc)
	} else {
		assert.NoError(t, err, "TEST: Unexpected Error: %s", tc.desc)
		assert.IsType(t, tc.expected, result, "TEST: %s: unexpected result type", tc.desc)
	}
}



