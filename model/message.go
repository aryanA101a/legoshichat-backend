package model

import "time"

type SendMessageByIDRequest struct {
	Content     string `json:"content" validate:"required"`
	RecipientID string `json:"recipientId" validate:"required"`
}
type SendMessageByPhoneNumberRequest struct {
	Content     string `json:"content" validate:"required"`
	RecipientPhoneNumber uint64 `json:"recipientPhoneNumber" validate:"required,min=1000000000,max=9999999999"`
}

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Timestamp time.Time `json:"timestamp"`
}
