package model

import "time"

type Profile struct {
	AllId      string
	Birthday   time.Time
	Firstname  string
	Lastname   string
	MiddleName string
	Faculty    string
	Grade      string
	Group      string
	Status     string
}

type CreateProfileInput struct {
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
