package query

import (
	"context"
)

type UserQuery interface {
	FindByEmail(ctx context.Context, email string) <-chan Result
	FindByEmailAndPassword(ctx context.Context, email, password string) <-chan Result
}

type TokenQuery interface {
	Find(key string) <-chan Result
}

type Result struct {
	Result interface{}
	Error  error
}
