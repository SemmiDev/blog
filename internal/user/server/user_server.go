package server

import (
	"github.com/SemmiDev/blog/internal/user/command"
	commandPostgresql "github.com/SemmiDev/blog/internal/user/command/postgresql"
	"github.com/SemmiDev/blog/internal/user/domain"
	"github.com/SemmiDev/blog/internal/user/domain/service"
	"github.com/SemmiDev/blog/internal/user/query"
	queryPostgresql "github.com/SemmiDev/blog/internal/user/query/postgresql"
	"github.com/SemmiDev/blog/internal/user/token"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type UserServer struct {
	UserCommand command.UserCommand
	UserQuery   query.UserQuery
	UserService domain.UserService
	TokenMaker  token.Maker
}

func NewUserServer(db *pgxpool.Pool, tokenMaker token.Maker) (*UserServer, error) {
	userServer := &UserServer{
		UserCommand: commandPostgresql.NewUserCommandPostgresql(db),
		UserQuery:   queryPostgresql.NewUserQueryPostgresql(db),
		TokenMaker:  tokenMaker,
	}
	userServer.UserService = service.UserServiceImpl{UserQuery: userServer.UserQuery}
	return userServer, nil
}

func (s *UserServer) Mount(r fiber.Router) {
	r.Use(s.AuthMiddleware())
	r.Post("/password/change", s.ChangePassword)
}

type ChangePasswordReq struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	NewConfirmPassword string `json:"new_confirm_password"`
}

func (s *UserServer) ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordReq
	if err := c.BodyParser(&req); err != nil {
		return Error(c, NewRequestValidationError(ParseFailed, "body"))
	}

	payload := c.Context().UserValue(authorizationPayloadKey).(*token.Payload)

	user, err := s.UserService.FindUserByEmail(c.Context(), payload.Email)
	if err != nil {
		return Error(c, err)
	}

	err = user.ChangePassword(req.OldPassword, req.NewPassword, req.NewConfirmPassword)
	if err != nil {
		return Error(c, err)
	}

	err = <-s.UserCommand.UpdatePassword(c.Context(), &user)
	if err != nil {
		return Error(c, err)
	}

	return c.SendStatus(http.StatusOK)
}
