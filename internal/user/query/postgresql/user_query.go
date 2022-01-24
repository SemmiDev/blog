package postgresql

import (
	"context"
	"errors"
	"github.com/SemmiDev/blog/internal/user/query"
	"github.com/SemmiDev/blog/internal/user/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserQueryPostgresql struct {
	DB *pgxpool.Pool
}

func NewUserQueryPostgresql(DB *pgxpool.Pool) *UserQueryPostgresql {
	return &UserQueryPostgresql{DB: DB}
}

type userReadResult struct {
	ID          string
	Name        string
	Nickname    string
	Email       string
	Password    []byte
	Bio         string
	Image       string
	CreatedDate time.Time
}

func (u UserQueryPostgresql) FindByEmail(ctx context.Context, email string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		userRead := storage.User{}
		rowsData := userReadResult{}

		err := u.DB.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(
			&rowsData.ID,
			&rowsData.Name,
			&rowsData.Nickname,
			&rowsData.Email,
			&rowsData.Password,
			&rowsData.Bio,
			&rowsData.Image,
			&rowsData.CreatedDate,
		)

		if err != nil {
			exists := strings.Contains(err.Error(), "no")
			if exists {
				result <- query.Result{Error: errors.New("account not found")}
			}
			result <- query.Result{Error: errors.New("internal server error")}
		}

		userRead = storage.User{
			ID:          rowsData.ID,
			Name:        rowsData.Name,
			Nickname:    rowsData.Nickname,
			Email:       rowsData.Email,
			Password:    rowsData.Password,
			Bio:         rowsData.Bio,
			Image:       rowsData.Image,
			CreatedDate: rowsData.CreatedDate,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}

func (u UserQueryPostgresql) FindByEmailAndPassword(ctx context.Context, email, password string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		userRead := storage.User{}
		rowsData := userReadResult{}

		err := u.DB.QueryRow(ctx, `SELECT * FROM users
			WHERE email = $1`, email).Scan(
			&rowsData.ID,
			&rowsData.Name,
			&rowsData.Nickname,
			&rowsData.Email,
			&rowsData.Password,
			&rowsData.Bio,
			&rowsData.Image,
			&rowsData.CreatedDate,
		)

		if err != nil {
			exists := strings.Contains(err.Error(), "no")
			if exists {
				result <- query.Result{Error: errors.New("account not found")}
			}
			result <- query.Result{Error: errors.New("internal server error")}
		}

		err = bcrypt.CompareHashAndPassword(rowsData.Password, []byte(password))
		if err != nil {
			result <- query.Result{Error: errors.New("incorrect password")}
		}

		userRead = storage.User{
			ID:          rowsData.ID,
			Name:        rowsData.Name,
			Nickname:    rowsData.Nickname,
			Email:       rowsData.Email,
			Password:    rowsData.Password,
			Bio:         rowsData.Bio,
			Image:       rowsData.Image,
			CreatedDate: rowsData.CreatedDate,
		}

		result <- query.Result{Result: userRead}
		close(result)
	}()

	return result
}
