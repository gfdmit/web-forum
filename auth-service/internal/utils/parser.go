package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gfdmit/web-forum/auth-service/internal/model"
)

const baseURL = "https://newlk.kpfu.ru/user-api"

func post(url string, payload any) ([]byte, error) {
	data, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func get(url, token string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func ParseKFU(login, password string) (*model.Profile, error) {
	loginBody, err := post(baseURL+"/login", map[string]string{
		"username": login,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	loginResp := model.LoginResponse{}
	json.Unmarshal(loginBody, &loginResp)

	switcherBody, err := get(baseURL+"/switcher", loginResp.Token)
	if err != nil {
		return nil, err
	}
	switcher := model.Switcher{}
	json.Unmarshal(switcherBody, &switcher)

	if len(switcher.Roles) == 0 {
		return nil, fmt.Errorf("no roles found for user")
	}

	birthday, err := time.Parse("02.01.2006", loginResp.User.Birthday)
	if err != nil {
		return nil, err
	}

	result := &model.Profile{
		AllId:      loginResp.User.AllID,
		Birthday:   birthday,
		Firstname:  loginResp.User.Firstname,
		Lastname:   loginResp.User.Lastname,
		MiddleName: loginResp.User.MiddleName,
		Faculty:    switcher.Roles[0].Faculty,
		Grade:      switcher.Roles[0].Grade,
		Group:      switcher.Roles[0].Group,
		Status:     switcher.Roles[0].Status,
	}

	return result, nil
}
