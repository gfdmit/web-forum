package model

import "time"

type Profile struct {
	ID           int
	UserID       int
	UniversityID string
	Firstname    string
	Lastname     string
	Middlename   string
	Birthday     time.Time
	Faculty      string
	Grade        string
	Group        string
	Status       string
}
