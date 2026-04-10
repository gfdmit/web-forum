package repository

import "time"

type Board struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Post struct {
	ID        int        `json:"id"`
	UserID    *int       `json:"user_id,omitempty"`
	BoardID   int        `json:"board_id"`
	Title     string     `json:"title,omitempty"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Comment struct {
	ID        int        `json:"id"`
	UserID    *int       `json:"user_id,omitempty"`
	PostID    int        `json:"post_id"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
