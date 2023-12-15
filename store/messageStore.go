package store

import (
	"fmt"

	"github.com/aryanA101a/legoshichat-backend/model"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/gofr"
)

type message struct {
}

type MessageStore interface {
	AddMessage(ctx *gofr.Context, message model.Message) error
}

func NewMessageStore(db *datastore.SQLClient) MessageStore {
	m := message{}
	err:=m.init(db)
	if err!=nil{
		fmt.Println("errr:",err)
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
