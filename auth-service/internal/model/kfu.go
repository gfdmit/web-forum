package model

type LoginResponse struct {
	Token string `json:"accessToken"`
	User  Data   `json:"user"`
}

type Data struct {
	AllID      string `json:"allId"`
	Birthday   string `json:"birthday"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	MiddleName string `json:"middleName"`
	Roles      []Role `json:"roles"`
}

type Role struct {
	UserID string `json:"userId"`
	Type   string `json:"type"`
}

type Switcher struct {
	Roles []RoleInfo `json:"roleinfo"`
}

type RoleInfo struct {
	Faculty string `json:"faculty"`
	Grade   string `json:"grade"`
	Group   string `json:"group"`
	Status  string `json:"status"`
}
