package storage

import (
	"time"
)

type User struct {
	ID          string
	Name        string
	Email       string
	Nickname    string
	Password    []byte
	Bio         string
	Image       string
	CreatedDate time.Time
	LastUpdated time.Time
}

type UserAuth struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	NickName    string `json:"nickname"`
	AccessToken string `json:"access_token"`
}
