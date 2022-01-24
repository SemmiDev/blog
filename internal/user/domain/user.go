package domain

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

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

type UserService interface {
	FindUserByEmail(ctx context.Context, email string) (User, error)
}

func CreateUser(email, name, password string) (*User, error) {
	if name == "" {
		return nil, UserError{UserErrorNameEmptyCode}
	}
	if email == "" {
		return nil, UserError{UserErrorEmailEmptyCode}
	}
	if password == "" {
		return nil, UserError{UserErrorPasswordEmptyCode}
	}
	if len(password) < 6 {
		return nil, UserError{UserErrorInvalidPasswordLengthCode}
	}
	if !mailRegex.MatchString(email) {
		return nil, UserError{UserErrorNameEmptyCode}
	}

	nickname := strings.Split(email, "@")[0]
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return &User{
		ID:       uuid.NewString(),
		Name:     name,
		Nickname: nickname,
		Email:    email,
		Password: hash,
	}, nil
}

func (u *User) IsPasswordValid(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false, UserError{UserErrorWrongPasswordCode}
	}

	return true, nil
}

func (u *User) ChangePassword(oldPassword, newPassword, newConfirmPassword string) error {
	valid, err := u.IsPasswordValid(oldPassword)
	if !valid || err != nil {
		return UserError{UserChangePasswordErrorWrongOldPasswordCode}
	}

	if newPassword == "" || newConfirmPassword == "" {
		return UserError{UserErrorPasswordEmptyCode}
	}
	if len(newPassword) < 6 || len(newConfirmPassword) < 6 {
		return UserError{UserErrorInvalidPasswordLengthCode}
	}
	if newPassword != newConfirmPassword {
		return UserError{UserErrorPasswordConfirmationNotMatchCode}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func (u *User) ResetPassword(newPassword, newConfirmPassword string) error {
	if newPassword == "" || newConfirmPassword == "" {
		return UserError{UserErrorPasswordEmptyCode}
	}
	if len(newPassword) < 6 || len(newConfirmPassword) < 6 {
		return UserError{UserErrorInvalidPasswordLengthCode}
	}
	if newPassword != newConfirmPassword {
		return UserError{UserErrorPasswordConfirmationNotMatchCode}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}
