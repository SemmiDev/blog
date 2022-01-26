package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// User represents a user table in the database.
type User struct {
	ID          string
	Name        string
	Nickname    string
	Email       string
	Password    []byte
	Bio         string
	Image       string
	CreatedDate time.Time
}

// CreateUser creates a new user and returns it.
func CreateUser(email, name, password string) (*User, error) {
	nickname := strings.Split(email, "@")[0]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:       uuid.NewString(),
		Name:     name,
		Nickname: nickname,
		Email:    email,
		Password: hash,
	}

	return &user, nil
}

// ChangePassword changes the user's password.
func (u *User) ChangePassword(oldPassword, newPassword, newConfirmPassword string) error {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(oldPassword))
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

// ResetPassword resets the user's password.
func (u *User) ResetPassword(newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}
