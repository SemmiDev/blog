package storage

import (
	"time"
)

// User is will be used as response for get details of user.
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

// UserAuth it will be used as response for authentication.
type UserAuth struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	NickName    string `json:"nickname"`
	AccessToken string `json:"access_token"`
}
