package store

import (
	"context"
	"fmt"

	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/aryanA101a/legoshichat-backend/model"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
)
func TestAddMessage(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	messageStore := message{}

	sampleMessage := model.Message{
		ID:        "test-message-id",
		Content:   "Hello, world!",
		From:      "sender-user-id",
		To:        "receiver-user-id",
		Timestamp: time.Now(),
	}

	mock.ExpectExec("INSERT INTO messages").
		WithArgs(sampleMessage.ID, sampleMessage.Content, sampleMessage.From, sampleMessage.To, sampleMessage.Timestamp).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	err = messageStore.AddMessage(ctx, sampleMessage)

	assert.NoError(t, err, "Unexpected error during message insertion")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("INSERT INTO messages").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf(""))

	err = messageStore.AddMessage(ctx, model.Message{})

	assert.Error(t, err, "Expected an error during failed message insertion")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
	
}

func TestGetMessage(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	messageStore := message{}

	userID := "test-user-id"
	messageID := "test-message-id"

	sampleMessage := model.Message{
		ID:        messageID,
		Content:   "Hello, world!",
		From:      userID,
		To:        "receiver-user-id",
		Timestamp: time.Now(),
	}

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "senderId", "recieverId", "timestamp"}).
			AddRow(sampleMessage.ID, sampleMessage.Content, sampleMessage.From, sampleMessage.To, sampleMessage.Timestamp))

	retrievedMessage, err := messageStore.GetMessage(ctx, userID, messageID)

	assert.NoError(t, err, "Unexpected error during message retrieval")
	assert.NotNil(t, retrievedMessage, "Expected a non-nil retrieved message")
	assert.Equal(t, messageID, retrievedMessage.ID, "Mismatch in retrieved message ID")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnError(fmt.Errorf(""))

	_, err = messageStore.GetMessage(ctx, userID, messageID)

	assert.Error(t, err, "Expected an error during failed message retrieval")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "senderId", "recieverId", "timestamp"}).
			AddRow(sampleMessage.ID, sampleMessage.Content, "different-user-id", sampleMessage.To, sampleMessage.Timestamp))

	_, err = messageStore.GetMessage(ctx, userID, messageID)

	assert.EqualError(t, err, e.NewError("You are not authorized to see that message").Error(), "Expected an authorization error")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdateMessage(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	messageStore := message{}

	userID := "test-user-id"
	messageID := "test-message-id"
	updatedContent := "Updated content"

	sampleMessage := model.Message{
		ID:        messageID,
		Content:   "Hello, world!",
		From:      userID,
		To:        "receiver-user-id",
		Timestamp: time.Now(),
	}

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "senderId", "recieverId", "timestamp"}).
			AddRow(sampleMessage.ID, sampleMessage.Content, sampleMessage.From, sampleMessage.To, sampleMessage.Timestamp))

	mock.ExpectExec("UPDATE messages SET content=").
		WithArgs(updatedContent, messageID).
		WillReturnResult(sqlmock.NewResult(0, 1)).
		WillReturnError(nil)

	updatedMessage, err := messageStore.UpdateMessage(ctx, userID, messageID, updatedContent)

	assert.NoError(t, err, "Unexpected error during message update")
	assert.NotNil(t, updatedMessage, "Expected a non-nil updated message")
	assert.Equal(t, updatedContent, updatedMessage.Content, "Mismatch in updated message content")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnError(fmt.Errorf(""))

	_, err = messageStore.UpdateMessage(ctx, userID, messageID, updatedContent)

	assert.Error(t, err, "Expected an error during failed message update")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "senderId", "recieverId", "timestamp"}).
			AddRow(sampleMessage.ID, sampleMessage.Content, "different-user-id", sampleMessage.To, sampleMessage.Timestamp))

	_, err = messageStore.UpdateMessage(ctx, userID, messageID, updatedContent)

	assert.EqualError(t, err, e.NewError("You are not authorized to update that message").Error(), "Expected an authorization error")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeleteMessage(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	messageStore := message{}

	userID := "test-user-id"
	messageID := "test-message-id"

	mock.ExpectQuery("SELECT senderId FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnRows(sqlmock.NewRows([]string{"senderId"}).
			AddRow(userID))

	mock.ExpectExec("DELETE FROM messages where id=").
		WithArgs(messageID).
		WillReturnResult(sqlmock.NewResult(0, 1)).
		WillReturnError(nil)

	err = messageStore.DeleteMessage(ctx, userID, messageID)

	assert.NoError(t, err, "Unexpected error during message deletion")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT senderId FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnError(fmt.Errorf(""))

	err = messageStore.DeleteMessage(ctx, userID, messageID)

	assert.Error(t, err, "Expected an error during failed message deletion")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT senderId FROM messages WHERE id=").
		WithArgs(messageID).
		WillReturnRows(sqlmock.NewRows([]string{"senderId"}).
			AddRow("different-user-id"))

	err = messageStore.DeleteMessage(ctx, userID, messageID)

	assert.EqualError(t, err, e.NewError("You are not authorized to delete that message").Error(), "Expected an authorization error")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetMessages(t *testing.T) {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %v", err)
	}
	defer db.Close()

	ctx.Context = context.Background()
	ctx.DataStore = datastore.DataStore{ORM: db}

	messageStore := message{}

	senderID := "test-sender-id"
	receiverID := "test-receiver-id"

	page := uint(1)
	limit := uint(10)

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages").
		WithArgs(senderID, receiverID, limit, (page-1)*limit).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "senderId", "recieverId", "timestamp"}).
			AddRow("message-id-1", "Hello", senderID, receiverID, time.Now()).
			AddRow("message-id-2", "Hi", senderID, receiverID, time.Now()))

	messages, err := messageStore.GetMessages(ctx, senderID, receiverID, page, limit)

	assert.NoError(t, err, "Unexpected error during message retrieval")
	assert.NotNil(t, messages, "Expected a non-nil list of messages")
	assert.Len(t, *messages, 2, "Unexpected number of retrieved messages")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id,content,senderId,recieverId,timestamp FROM messages").
		WithArgs(senderID, receiverID, limit, (page-1)*limit).
		WillReturnError(fmt.Errorf(""))

	_, err = messageStore.GetMessages(ctx, senderID, receiverID, page, limit)

	assert.Error(t, err, "Expected an error during failed message retrieval")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}