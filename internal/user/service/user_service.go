package service

import (
	cloud "cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/SemmiDev/blog/config"
	"github.com/SemmiDev/blog/internal/common/random"
	"github.com/SemmiDev/blog/internal/user/entity"
	. "github.com/SemmiDev/blog/internal/user/helper"
	"github.com/SemmiDev/blog/internal/user/query"
	"github.com/SemmiDev/blog/internal/user/repository"
	"github.com/SemmiDev/blog/internal/user/storage"
	"github.com/SemmiDev/blog/internal/user/token"
	"log"
	"mime/multipart"
	"strings"
	"time"
)

// UserServiceImpl is a struct that implements UserService interface.
type UserServiceImpl struct {
	UserQuery    query.UserQuery
	TokenQuery   query.TokenQuery
	UserCommand  repository.UserCommand
	TokenCommand repository.TokenCommand
	TokenMaker   token.Maker
	CloudStorage *cloud.Client
}

// SendVerificationCode sends verification code to user's email.
func (s *UserServiceImpl) SendVerificationCode(ctx context.Context, email string, kind string) error {
	if email == "" {
		return NewErr(ErrEmailEmptyCode, "email")
	}
	if !MailRegex.MatchString(email) {
		return NewErr(ErrInvalidEmailCode, "email")
	}

	key := random.Codes(10)
	val := fmt.Sprintf("%s|%s|%s", key, email, kind)

	// log the value for now
	log.Println(val)

	err := <-s.TokenCommand.Set(key, []byte(val), 30*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

// RegisterNewUser registers new user.
func (s *UserServiceImpl) RegisterNewUser(ctx context.Context, code, name, password string) (storage.UserAuth, error) {
	if len(code) != 10 {
		return storage.UserAuth{}, NewErr(ErrInvalidCode, "code")
	}
	if !NumberRegex.MatchString(code) {
		return storage.UserAuth{}, NewErr(ErrInvalidCode, "code")
	}
	if name == "" {
		return storage.UserAuth{}, NewErr(ErrNameEmptyCode, "name")
	}
	if password == "" {
		return storage.UserAuth{}, NewErr(ErrPasswordEmptyCode, "password")
	}
	if len(password) < 6 {
		return storage.UserAuth{}, NewErr(ErrInvalidPasswordLengthCode, "password")
	}

	tokenResult := <-s.TokenQuery.Find(code)
	if tokenResult.Error != nil {
		return storage.UserAuth{}, NewErr(ErrInvalidCode, "code")
	}

	codeVerification, _ := tokenResult.Result.([]byte)
	if codeVerification == nil {
		return storage.UserAuth{}, NewErr(ErrInvalidCode, "code")
	}

	extractCode := strings.Split(string(codeVerification), "|")
	if extractCode[0] != code {
		return storage.UserAuth{}, NewErr(ErrInvalidCode, "code")
	}
	email := extractCode[1]

	err := <-s.TokenCommand.Delete(code)
	if err != nil {
		return storage.UserAuth{}, err
	}

	user, err := entity.CreateUser(email, name, password)
	if err != nil {
		return storage.UserAuth{}, err
	}

	err = <-s.UserCommand.Save(ctx, user)
	if err != nil {
		return storage.UserAuth{}, err
	}

	accessToken, err := s.TokenMaker.CreateToken(user.Email, config.Env.AccessTokenDuration)
	if err != nil {
		return storage.UserAuth{}, err
	}

	auth := storage.UserAuth{
		UserID:      user.ID,
		Name:        user.Name,
		Email:       user.Email,
		NickName:    user.Nickname,
		AccessToken: accessToken,
	}

	return auth, nil
}

// Authorize user by email and password.
func (s *UserServiceImpl) Authorize(ctx context.Context, email string, password string) (storage.UserAuth, error) {
	if email == "" {
		return storage.UserAuth{}, NewErr(ErrEmailEmptyCode, "email")
	}
	if !MailRegex.MatchString(email) {
		return storage.UserAuth{}, NewErr(ErrInvalidEmailCode, "email")
	}
	if password == "" {
		return storage.UserAuth{}, NewErr(ErrPasswordEmptyCode, "password")
	}
	if len(password) < 6 {
		return storage.UserAuth{}, NewErr(ErrInvalidPasswordLengthCode, "password")
	}

	queryResult := <-s.UserQuery.FindByEmailAndPassword(ctx, email, password)
	if queryResult.Error != nil {
		return storage.UserAuth{}, queryResult.Error
	}

	user, ok := queryResult.Result.(storage.User)
	if !ok {
		return storage.UserAuth{}, errors.New("helper type assertion")
	}

	accessToken, _ := s.TokenMaker.CreateToken(user.Email, config.Env.AccessTokenDuration)
	userAuth := storage.UserAuth{
		UserID:      user.ID,
		Name:        user.Name,
		Email:       user.Email,
		NickName:    user.Nickname,
		AccessToken: accessToken,
	}
	return userAuth, nil
}

// ResetPassword resets user's password.
func (s *UserServiceImpl) ResetPassword(ctx context.Context, code, newPassword, newConfirmPassword string) error {
	if len(code) != 10 {
		return NewErr(ErrInvalidCode, "code")
	}
	if !NumberRegex.MatchString(code) {
		return NewErr(ErrInvalidCode, "code")
	}
	if newPassword == "" {
		return NewErr(ErrPasswordEmptyCode, "new password")
	}
	if newConfirmPassword == "" {
		return NewErr(ErrPasswordEmptyCode, "new confirm password")
	}
	if len(newPassword) < 6 {
		return NewErr(ErrInvalidPasswordLengthCode, "new password")
	}
	if len(newConfirmPassword) < 6 {
		return NewErr(ErrInvalidPasswordLengthCode, "new confirm password")
	}
	if newPassword != newConfirmPassword {
		return NewErr(ErrPasswordConfirmationNotMatchCode, "password")
	}

	// business logic
	tokenResult := <-s.TokenQuery.Find(code)
	if tokenResult.Error != nil {
		return NewErr(ErrInvalidCode, "code")
	}

	codeVerification, _ := tokenResult.Result.([]byte)
	if codeVerification == nil {
		return NewErr(ErrInvalidCode, "code")
	}

	extract := strings.Split(string(codeVerification), "|")
	if extract[0] != code {
		return NewErr(ErrInvalidCode, "code")
	}

	user, err := s.FindUserByEmail(ctx, extract[1])
	if err != nil {
		return err
	}

	err = user.ResetPassword(newPassword)
	if err != nil {
		return err
	}

	err = <-s.UserCommand.UpdatePassword(ctx, &user)
	if err != nil {
		return err
	}

	err = <-s.TokenCommand.Delete(code)
	if err != nil {
		return err
	}

	return nil
}

// ChangePassword changes user's password.
func (s *UserServiceImpl) ChangePassword(ctx context.Context, oldPassword, newPassword, newConfirmPassword string, email string) error {
	if newPassword == "" {
		return NewErr(ErrPasswordEmptyCode, "new password")
	}
	if newConfirmPassword == "" {
		return NewErr(ErrPasswordEmptyCode, "new confirm password")
	}
	if len(newPassword) < 6 {
		return NewErr(ErrInvalidPasswordLengthCode, "new password")
	}
	if len(newConfirmPassword) < 6 {
		return NewErr(ErrInvalidPasswordLengthCode, "new confirm password")
	}
	if newPassword != newConfirmPassword {
		return NewErr(ErrPasswordConfirmationNotMatchCode, "password")
	}

	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = user.ChangePassword(oldPassword, newPassword, newConfirmPassword)
	if err != nil {
		return err
	}

	err = <-s.UserCommand.UpdatePassword(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

// ChangeBio changes user's bio.
func (s *UserServiceImpl) ChangeBio(ctx context.Context, bio, email string) error {
	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	user.Bio = bio
	err = <-s.UserCommand.UpdateBio(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

// ChangeImage changes user's image.
func (s *UserServiceImpl) ChangeImage(ctx context.Context, file *multipart.FileHeader, email string) error {
	result := <-UploadImage(ctx, s.CloudStorage, file, email)
	if result.Error != nil {
		return result.Error
	}

	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	user.Image = result.Path
	err = <-s.UserCommand.UpdateImage(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
