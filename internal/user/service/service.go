package service

import (
	"context"
	"errors"
	"github.com/SemmiDev/blog/internal/user/entity"
	"github.com/SemmiDev/blog/internal/user/storage"
	"mime/multipart"
)

// UserService is a service for managing users.
type UserService interface {
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	SendVerificationCode(ctx context.Context, email string, kind string) error
	RegisterNewUser(ctx context.Context, code, name, password string) (storage.UserAuth, error)
	Authorize(ctx context.Context, email, password string) (storage.UserAuth, error)
	ResetPassword(ctx context.Context, code, newPassword, newConfirmPassword string) error
	ChangePassword(ctx context.Context, oldPassword, newPassword, newConfirmPassword, email string) error
	ChangeBio(ctx context.Context, bio, email string) error
	ChangeImage(ctx context.Context, file *multipart.FileHeader, email string) error
}

// FindUserByEmail returns a user by email.
func (s UserServiceImpl) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	result := <-s.UserQuery.FindByEmail(ctx, email)

	if result.Error != nil {
		return entity.User{}, result.Error
	}

	user, ok := result.Result.(storage.User)
	if !ok {
		return entity.User{}, errors.New("helper type assertion")
	}

	userResult := entity.User{
		ID:          user.ID,
		Name:        user.Name,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Password:    user.Password,
		Bio:         user.Bio,
		Image:       user.Image,
		CreatedDate: user.CreatedDate,
	}

	return userResult, nil
}
