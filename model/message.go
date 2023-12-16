package model

import "time"

const RequestMessageLimit=5

type SendMessageByIDRequest struct {
	Content     string `json:"content" validate:"required,min=1"`
	RecipientID string `json:"recipientId" validate:"required"`
}
type SendMessageByPhoneNumberRequest struct {
	Content     string `json:"content" validate:"required,min=1"`
	RecipientPhoneNumber uint64 `json:"recipientPhoneNumber" validate:"required,min=1000000000,max=9999999999"`
}
type UpdateMessageRequest struct {
	Content     string `json:"content" validate:"required,min=1"`
}
type GetMessagesRequest struct {
	Page uint `json:"page" validate:"required"`
	SenderID string `json:"senderId" validate:"required"`
	RecipientID string `json:"recipientId" validate:"required"`
}

type GetMessagesResponse struct {
	Page uint `json:"page" validate:"required"`
	LastPage bool `json:"lastPage" validate:"required"`
	Messages []Message `json:"messages" validate:"required"`
}

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Timestamp time.Time `json:"timestamp"`
}
