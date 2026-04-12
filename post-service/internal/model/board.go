package model

import "time"

type Board struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type CreateBoardInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
