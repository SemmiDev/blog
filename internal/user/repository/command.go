package repository

import (
	"context"
	"github.com/SemmiDev/blog/internal/user/entity"
	"time"
)

type UserCommand interface {
	UserSaver
	UserUpdater
}

type UserSaver interface {
	Save(ctx context.Context, arg *entity.User) <-chan error
}

type UserUpdater interface {
	UpdatePassword(ctx context.Context, arg *entity.User) <-chan error
	UpdateBio(ctx context.Context, arg *entity.User) <-chan error
	UpdateImage(ctx context.Context, arg *entity.User) <-chan error
}

type TokenCommand interface {
	Set(key string, val []byte, exp time.Duration) <-chan error
	Delete(key string) <-chan error
}

type Result struct {
	Result interface{}
	Error  error
}
