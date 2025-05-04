package models

import "time"

type ChatMessage struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessageResponse struct {
	ID        int64        `json:"id"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}
