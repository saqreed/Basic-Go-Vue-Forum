package models

import (
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	PostID    int64     `json:"post_id"`
	AuthorID  int64     `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentResponse struct {
	ID        int64        `json:"id"`
	Content   string       `json:"content"`
	PostID    int64        `json:"post_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
