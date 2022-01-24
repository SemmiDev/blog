package service

import (
	"context"
	"errors"
	"github.com/SemmiDev/blog/internal/user/domain"
	"github.com/SemmiDev/blog/internal/user/query"
	"github.com/SemmiDev/blog/internal/user/storage"
)

type UserServiceImpl struct {
	UserQuery query.UserQuery
}

func (s UserServiceImpl) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	result := <-s.UserQuery.FindByEmail(ctx, email)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	user, ok := result.Result.(storage.User)
	if !ok {
		return domain.User{}, errors.New("error type assertion")
	}

	return domain.User{
		ID:          user.ID,
		Name:        user.Name,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Password:    user.Password,
		Bio:         user.Bio,
		Image:       user.Image,
		CreatedDate: user.CreatedDate,
	}, nil
}
