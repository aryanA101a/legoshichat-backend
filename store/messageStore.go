package store

import (
	"fmt"

	e "github.com/aryanA101a/legoshichat-backend/error"
	"github.com/aryanA101a/legoshichat-backend/model"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
)

type message struct {
}

type MessageStore interface {
	AddMessage(ctx *gofr.Context, message model.Message) error
	GetMessage(ctx *gofr.Context, userId, messageId string) (*model.Message, error)
	UpdateMessage(ctx *gofr.Context, userId, messageId, updatedContent string) (*model.Message, error)
	DeleteMessage(ctx *gofr.Context, userId, messageId string) error
	GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error)
}

func NewMessageStore(db *datastore.SQLClient) MessageStore {
	m := message{}
	err := m.init(db)
	if err != nil {
		fmt.Println("errr:", err)
	}
	return m
}

func (m message) init(db *datastore.SQLClient) error {
	return m.createMessageTable(db)
}

func (m message) AddMessage(ctx *gofr.Context, message model.Message) error {
	_, err := ctx.DB().ExecContext(ctx, "INSERT INTO messages (id,content,senderId,recieverId,timestamp) VALUES ($1,$2,$3,$4,$5)", message.ID, message.Content, message.From, message.To, message.Timestamp)
	return err
}

func (m message) GetMessage(ctx *gofr.Context, userId, messageId string) (*model.Message, error) {
	var message model.Message
	err := ctx.DB().QueryRowContext(ctx, "SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=$1", messageId).
		Scan(&message.ID, &message.Content, &message.From, &message.To, &message.Timestamp)

	if err != nil {
		return nil, err
	}
	if !(userId == message.From || userId == message.To) {
		return nil, e.NewError("You are not authorized to see that message")
	}

	return &message, nil
}

func (m message) UpdateMessage(ctx *gofr.Context, userId, messageId, updatedContent string) (*model.Message, error) {
	var message model.Message
	err := ctx.DB().QueryRowContext(ctx, "SELECT id,content,senderId,recieverId,timestamp FROM messages WHERE id=$1", messageId).
		Scan(&message.ID, &message.Content, &message.From, &message.To, &message.Timestamp)

	if err != nil {
		return nil, err
	}
	if !(userId == message.From) {
		return nil, e.NewError("You are not authorized to update that message")
	}

	_, err = ctx.DB().ExecContext(ctx, "UPDATE messages SET content=$1 WHERE id=$2", updatedContent, messageId)
	if err != nil {
		return nil, err
	}

	message.Content = updatedContent
	return &message, nil
}

func (m message) DeleteMessage(ctx *gofr.Context, userId, messageId string) error {
	var senderId string
	err := ctx.DB().QueryRowContext(ctx, "SELECT senderId FROM messages WHERE id=$1", messageId).
		Scan(&senderId)

	if err != nil {
		return err
	}
	if !(userId == senderId) {
		return e.NewError("You are not authorized to delete that message")
	}
	_, err = ctx.DB().ExecContext(ctx, "DELETE FROM messages where id=$1 ", messageId)
	return err
}

func (m message) GetMessages(ctx *gofr.Context, senderId, recieverId string, page, limit uint) (*[]model.Message, error) {
	
	query:=`SELECT id,content,senderId,recieverId,timestamp FROM messages
	WHERE senderId=$1 and recieverId=$2 ORDER BY timestamp DESC LIMIT $3 OFFSET $4`
	
	rows, err := ctx.DB().QueryContext(ctx, query, senderId, recieverId, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := make([]model.Message, 0)

	for rows.Next() {
		var message model.Message

		err = rows.Scan(&message.ID, &message.Content, &message.From, &message.To, &message.Timestamp)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &messages, nil
}

func (message) createMessageTable(db *datastore.SQLClient) error {
	query := `CREATE TABLE IF NOT EXISTS messages (
		id UUID PRIMARY KEY,
		content TEXT NOT NULL,
		senderId UUID NOT NULL,
		recieverId UUID NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}
