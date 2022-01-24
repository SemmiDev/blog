package server

import (
	"context"
	"errors"
	"fmt"
	// load all configs
	. "github.com/SemmiDev/blog/config"
	"github.com/SemmiDev/blog/internal/common/memory"
	"github.com/SemmiDev/blog/internal/common/random"
	"github.com/SemmiDev/blog/internal/user/command"
	commandMemory "github.com/SemmiDev/blog/internal/user/command/memory"
	commandPostgresql "github.com/SemmiDev/blog/internal/user/command/postgresql"
	"github.com/SemmiDev/blog/internal/user/domain"
	"github.com/SemmiDev/blog/internal/user/domain/service"
	"github.com/SemmiDev/blog/internal/user/query"
	queryMemory "github.com/SemmiDev/blog/internal/user/query/memory"
	queryPostgresql "github.com/SemmiDev/blog/internal/user/query/postgresql"
	"github.com/SemmiDev/blog/internal/user/storage"
	"github.com/SemmiDev/blog/internal/user/token"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthServer struct {
	UserCommand  command.UserCommand
	TokenCommand command.TokenCommand
	UserQuery    query.UserQuery
	TokenQuery   query.TokenQuery
	UserService  domain.UserService
	TokenMaker   token.Maker
}

func NewAuthServer(db *pgxpool.Pool, tokenMaker token.Maker) (*AuthServer, error) {
	m := memory.New()
	userServer := &AuthServer{
		UserCommand:  commandPostgresql.NewUserCommandPostgresql(db),
		TokenCommand: commandMemory.NewTokenCommandMemory(m),
		UserQuery:    queryPostgresql.NewUserQueryPostgresql(db),
		TokenQuery:   queryMemory.NewTokenQueryMemory(m),
		TokenMaker:   tokenMaker,
	}

	userServer.UserService = service.UserServiceImpl{UserQuery: userServer.UserQuery}
	return userServer, nil
}

func (s *AuthServer) Mount(r fiber.Router) {
	r.Post("/register", s.Register)
	r.Post("/authorize", s.Authorize)
	r.Post("/password/reset", s.Reset)
}

type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthServer) Register(c *fiber.Ctx) error {
	var req RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return Error(c, NewRequestValidationError(ParseFailed, "body"))
	}

	code := c.Query("code")
	if code != "" {
		if len(code) != 10 {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}
		if !numberRegex.MatchString(code) {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		tokenResult := <-s.TokenQuery.Find(code)
		if tokenResult.Error != nil {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		codeVerification, _ := tokenResult.Result.([]byte)
		if codeVerification == nil {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		extract := strings.Split(string(codeVerification), "|")
		if extract[0] != code {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		userAuth, err := s.RegisterNewUser(c.Context(), extract[1], req.Name, req.Password)
		if err != nil {
			return Error(c, err)
		}

		err = <-s.TokenCommand.Delete(code)
		if err != nil {
			return Error(c, err)
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"data": userAuth,
		})
	}

	if req.Email == "" {
		return Error(c, NewRequestValidationError(Required, "email"))
	}
	if !mailRegex.MatchString(req.Email) {
		return Error(c, NewRequestValidationError(Invalid, "email"))
	}

	key := random.Codes(10)
	value := fmt.Sprintf("%s|%s", key, req.Email)

	log.Println(value)

	err := <-s.TokenCommand.Set(key, []byte(value), 30*time.Minute)
	if err != nil {
		return Error(c, err)
	}

	return c.SendStatus(http.StatusOK)
}

func (s *AuthServer) RegisterNewUser(ctx context.Context, email, name, password string) (*storage.UserAuth, error) {
	user, err := domain.CreateUser(email, name, password)
	if err != nil {
		return nil, err
	}

	err = <-s.UserCommand.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, _ := s.TokenMaker.CreateToken(user.Email, Config.AccessTokenDuration)

	return &storage.UserAuth{
		UserID:      user.ID,
		Name:        user.Name,
		Email:       user.Email,
		NickName:    user.Nickname,
		AccessToken: accessToken,
	}, nil
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthServer) Authorize(c *fiber.Ctx) error {
	var req LoginReq
	if err := c.BodyParser(&req); err != nil {
		return Error(c, NewRequestValidationError(ParseFailed, "body"))
	}

	queryResult := <-s.UserQuery.FindByEmailAndPassword(c.Context(), req.Email, req.Password)
	if queryResult.Error != nil {
		return Error(c, queryResult.Error)
	}

	user, ok := queryResult.Result.(storage.User)
	if !ok {
		return Error(c, errors.New("error type assertion"))
	}

	accessToken, _ := s.TokenMaker.CreateToken(user.Email, Config.AccessTokenDuration)
	userAuth := storage.UserAuth{
		UserID:      user.ID,
		Name:        user.Name,
		Email:       user.Email,
		NickName:    user.Nickname,
		AccessToken: accessToken,
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": userAuth,
	})
}

type ResetPasswordReq struct {
	Email              string `json:"email"`
	NewPassword        string `json:"new_password"`
	NewConfirmPassword string `json:"new_confirm_password"`
}

func (s *AuthServer) Reset(c *fiber.Ctx) error {
	var req ResetPasswordReq
	if err := c.BodyParser(&req); err != nil {
		return Error(c, NewRequestValidationError(ParseFailed, "body"))
	}

	code := c.Query("code")
	if code != "" {
		if len(code) != 10 {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}
		if !numberRegex.MatchString(code) {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		tokenResult := <-s.TokenQuery.Find(code)
		if tokenResult.Error != nil {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		codeVerification, _ := tokenResult.Result.([]byte)
		if codeVerification == nil {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		extract := strings.Split(string(codeVerification), "|")
		if extract[0] != code {
			return Error(c, NewRequestValidationError(Invalid, "code"))
		}

		user, err := s.UserService.FindUserByEmail(c.Context(), extract[1])
		if err != nil {
			return Error(c, err)
		}

		err = user.ResetPassword(req.NewPassword, req.NewConfirmPassword)
		if err != nil {
			return Error(c, err)
		}

		err = <-s.UserCommand.UpdatePassword(c.Context(), &user)
		if err != nil {
			return Error(c, err)
		}

		err = <-s.TokenCommand.Delete(code)
		if err != nil {
			return Error(c, err)
		}

		return c.SendStatus(http.StatusOK)
	}

	if req.Email == "" {
		return Error(c, NewRequestValidationError(Required, "email"))
	}
	if !mailRegex.MatchString(req.Email) {
		return Error(c, NewRequestValidationError(Invalid, "email"))
	}

	key := random.Codes(10)
	value := fmt.Sprintf("%s|%s", key, req.Email)

	log.Println(value)

	err := <-s.TokenCommand.Set(key, []byte(value), 30*time.Minute)
	if err != nil {
		return Error(c, err)
	}

	return c.SendStatus(http.StatusOK)
}
