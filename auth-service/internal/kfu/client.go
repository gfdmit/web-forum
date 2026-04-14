package kfu

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gfdmit/web-forum/auth-service/internal/model"
)

const baseURL = "https://newlk.kpfu.ru/user-api"

var ErrInvalidCredentials = errors.New("invalid credentials")

type KFUClient interface {
	ParseKFU(ctx context.Context, login, password string) (*model.Profile, error)
}

type client struct {
	http *http.Client
}

func NewClient() KFUClient {
	return &client{
		http: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *client) ParseKFU(ctx context.Context, login, password string) (*model.Profile, error) {
	loginResp, err := c.login(ctx, login, password)
	if err != nil {
		return nil, err
	}

	switcher, err := c.getSwitcher(ctx, loginResp.Token)
	if err != nil {
		return nil, err
	}

	if len(switcher.Roles) == 0 {
		return nil, fmt.Errorf("ParseKFU: %w", ErrInvalidCredentials)
	}

	birthday, err := time.Parse("02.01.2006", loginResp.User.Birthday)
	if err != nil {
		return nil, fmt.Errorf("ParseKFU parse birthday: %w", err)
	}

	return &model.Profile{
		AllId:      loginResp.User.AllID,
		Birthday:   birthday,
		Firstname:  loginResp.User.Firstname,
		Lastname:   loginResp.User.Lastname,
		MiddleName: loginResp.User.MiddleName,
		Faculty:    switcher.Roles[0].Faculty,
		Grade:      switcher.Roles[0].Grade,
		Group:      switcher.Roles[0].Group,
		Status:     switcher.Roles[0].Status,
	}, nil
}

func (c *client) login(ctx context.Context, login, password string) (model.LoginResponse, error) {
	body, err := json.Marshal(map[string]string{
		"username": login,
		"password": password,
	})
	if err != nil {
		return model.LoginResponse{}, fmt.Errorf("login marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/login", bytes.NewReader(body))
	if err != nil {
		return model.LoginResponse{}, fmt.Errorf("login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	data, err := c.do(req)
	if err != nil {
		return model.LoginResponse{}, err
	}

	var resp model.LoginResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return model.LoginResponse{}, fmt.Errorf("login unmarshal: %w", err)
	}
	return resp, nil
}

func (c *client) getSwitcher(ctx context.Context, token string) (model.Switcher, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/switcher", nil)
	if err != nil {
		return model.Switcher{}, fmt.Errorf("getSwitcher request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	data, err := c.do(req)
	if err != nil {
		return model.Switcher{}, err
	}

	var switcher model.Switcher
	if err := json.Unmarshal(data, &switcher); err != nil {
		return model.Switcher{}, fmt.Errorf("getSwitcher unmarshal: %w", err)
	}
	return switcher, nil
}

func (c *client) do(req *http.Request) ([]byte, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
