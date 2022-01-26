package postgresql

import (
	"context"
	"errors"
	"github.com/SemmiDev/blog/internal/user/entity"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserCommandPostgresql struct {
	DB *pgxpool.Pool
}

func NewUserCommandPostgresql(DB *pgxpool.Pool) *UserCommandPostgresql {
	return &UserCommandPostgresql{DB: DB}
}

func (u *UserCommandPostgresql) Save(ctx context.Context, arg *entity.User) <-chan error {
	result := make(chan error)

	go func() {
		count := 0
		err := u.DB.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE email = $1`, arg.Email).Scan(&count)
		if err != nil {
			result <- err
		}

		if count > 0 {
			result <- errors.New("account already exists")
		} else {
			_, err := u.DB.Exec(ctx, `INSERT INTO users (id, name, nickname, email, password) VALUES ($1, $2, $3, $4, $5)`,
				arg.ID, arg.Name, arg.Nickname, arg.Email, arg.Password)
			if err != nil {
				result <- err
			}
		}

		result <- nil
		close(result)
	}()

	return result
}

func (u *UserCommandPostgresql) UpdatePassword(ctx context.Context, arg *entity.User) <-chan error {
	result := make(chan error)

	go func() {
		_, err := u.DB.Exec(ctx, `UPDATE users SET password = $2 WHERE id = $1`,
			arg.ID, arg.Password)
		if err != nil {
			result <- err
		}

		result <- nil
		close(result)
	}()

	return result
}

func (u *UserCommandPostgresql) UpdateBio(ctx context.Context, arg *entity.User) <-chan error {
	result := make(chan error)

	go func() {
		_, err := u.DB.Exec(ctx, `UPDATE users SET bio = $2 WHERE id = $1`, arg.ID, arg.Bio)
		if err != nil {
			result <- err
		}

		result <- nil
		close(result)
	}()

	return result
}

func (u *UserCommandPostgresql) UpdateImage(ctx context.Context, arg *entity.User) <-chan error {
	result := make(chan error)

	go func() {
		_, err := u.DB.Exec(ctx, `UPDATE users SET image = $2 WHERE id = $1`, arg.ID, arg.Image)
		if err != nil {
			result <- err
		}

		result <- nil
		close(result)
	}()

	return result
}
