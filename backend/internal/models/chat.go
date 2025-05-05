package models

import "time"

type ChatMessage struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	ReplyToID *int      `json:"reply_to_id,omitempty"`
	ReplyTo   *ReplyTo  `json:"reply_to,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ReplyTo struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

type ChatMessageRequest struct {
	Content   string `json:"content"`
	ReplyToID *int   `json:"reply_to_id,omitempty"`
}

type ChatMessageResponse struct {
	ID        int64        `json:"id"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}
