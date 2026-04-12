package model

import "time"

type Comment struct {
	ID        int        `json:"id"`
	UserID    *int       `json:"user_id,omitempty"`
	PostID    int        `json:"post_id"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateCommentInput struct {
	UserID *int   `json:"-"`
	PostID int    `json:"post_id"`
	Text   string `json:"text"`
}
