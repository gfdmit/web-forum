package model

import "time"

type Profile struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	UniversityID string    `json:"university_id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Middlename   string    `json:"middlename"`
	Birthday     time.Time `json:"birthday"`
	Faculty      string    `json:"faculty"`
	Grade        string    `json:"grade"`
	Group        string    `json:"group"`
	Status       string    `json:"status"`
}
