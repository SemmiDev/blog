package command

import (
	"context"
	"github.com/SemmiDev/blog/internal/user/domain"
	"time"
)

type UserCommand interface {
	UserSaver
	UserUpdater
}

type UserSaver interface {
	Save(ctx context.Context, arg *domain.User) <-chan error
}

type UserUpdater interface {
	UpdatePassword(ctx context.Context, arg *domain.User) <-chan error
}

type TokenCommand interface {
	Set(key string, val []byte, exp time.Duration) <-chan error
	Delete(key string) <-chan error
}

type Result struct {
	Result interface{}
	Error  error
}
