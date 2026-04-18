package model

import "time"

type Post struct {
	ID        int        `json:"id"`
	UserID    *int       `json:"user_id,omitempty"`
	Author    *Author    `json:"authon,omitempty"`
	BoardID   int        `json:"board_id"`
	Title     string     `json:"title"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreatePostInput struct {
	UserID  *int   `json:"-"`
	BoardID int    `json:"board_id"`
	Title   string `json:"title"`
	Text    string `json:"text"`
}
