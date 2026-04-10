package model

import "time"

type Role struct {
	UserID string `json:"userId"`
	Type   string `json:"type"`
}

type User struct {
	AllID      string `json:"allId"`
	Birthday   string `json:"birthday"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	MiddleName string `json:"middleName"`
	Roles      []Role `json:"roles"`
}

type LoginResponse struct {
	Token string `json:"accessToken"`
	User  User   `json:"user"`
}

type RoleInfo struct {
	Faculty string `json:"faculty"`
	Grade   string `json:"grade"`
	Group   string `json:"group"`
	Status  string `json:"status"`
}

type Switcher struct {
	Roles []RoleInfo `json:"roleinfo"`
}

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
