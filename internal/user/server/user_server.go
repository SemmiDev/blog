package server

import (
	cloud "cloud.google.com/go/storage"
	"context"
	"github.com/SemmiDev/blog/config"
	zerolog "github.com/SemmiDev/blog/internal/common/logger"
	"github.com/SemmiDev/blog/internal/common/memory"
	"github.com/SemmiDev/blog/internal/user/helper"
	queryMemory "github.com/SemmiDev/blog/internal/user/query/memory"
	queryPostgresql "github.com/SemmiDev/blog/internal/user/query/postgresql"
	commandMemory "github.com/SemmiDev/blog/internal/user/repository/memory"
	commandPostgresql "github.com/SemmiDev/blog/internal/user/repository/postgresql"
	"github.com/SemmiDev/blog/internal/user/service"
	"github.com/SemmiDev/blog/internal/user/token"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/api/option"
	"net/http"
)

// UserServer is a struct that contains UserService for interacting with the service.
// and TokenMaker for generating & validating tokens.
type UserServer struct {
	UserService service.UserService
	TokenMaker  token.Maker
}

// NewUserServer returns a new UserServer.
func NewUserServer(db *pgxpool.Pool, tokenMaker token.Maker) (*UserServer, error) {
	// for now, we're using memory repository for stores code verification.
	// in the future, we'll use redis, or other stores.
	m := memory.New()

	// setup cloud client for interacting with google cloud storage.
	// it will be used for storing user avatar/profile images.
	opt := option.WithCredentialsFile(config.Env.FirebaseCredentialJSON)
	cloudClient, err := cloud.NewClient(context.Background(), opt)
	if err != nil {
		zerolog.Log.Error().Interface("cloud", err).Send()
		return nil, err
	}

	userServiceImpl := &service.UserServiceImpl{
		UserQuery:    queryPostgresql.NewUserQueryPostgresql(db),
		TokenQuery:   queryMemory.NewTokenQueryMemory(m),
		UserCommand:  commandPostgresql.NewUserCommandPostgresql(db),
		TokenCommand: commandMemory.NewTokenCommandMemory(m),
		TokenMaker:   tokenMaker,
		CloudStorage: cloudClient,
	}

	return &UserServer{
		UserService: userServiceImpl,
		TokenMaker:  tokenMaker,
	}, nil
}

// Mount mounts the UserServer to the fiber app.
func (s *UserServer) Mount(r fiber.Router) {
	// auth middleware for all routes.
	r.Use(s.AuthMiddleware())

	r.Put("/password/change", s.ChangePasswordHandler)
	r.Put("/profile/bio", s.ChangeBioHandler)
	r.Put("/profile/image", s.ChangeImageHandler)
}

// ChangePasswordHandler changes the user's password.
func (s *UserServer) ChangePasswordHandler(c *fiber.Ctx) error {
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")
	newConfirmPassword := c.FormValue("new_confirm_password")

	// get payload from context.
	payload := c.Context().UserValue(authorizationPayloadKey).(*token.Payload)

	err := s.UserService.ChangePassword(c.Context(), oldPassword, newPassword, newConfirmPassword, payload.Email)
	if err != nil {
		return helper.Error(c, err)
	}

	return c.SendStatus(http.StatusOK)
}

// ChangeBioHandler changes the user's bio.
func (s *UserServer) ChangeBioHandler(c *fiber.Ctx) error {
	bio := c.FormValue("bio")

	// get payload from context.
	payload := c.Context().UserValue(authorizationPayloadKey).(*token.Payload)

	err := s.UserService.ChangeBio(c.Context(), bio, payload.Email)
	if err != nil {
		return helper.Error(c, err)
	}

	return c.SendStatus(http.StatusOK)
}

// ChangeImageHandler changes the user's image.
func (s *UserServer) ChangeImageHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return helper.Error(c, helper.NewErr(helper.ErrParseCode, "image"))
	}

	// get payload from context.
	payload := c.Context().UserValue(authorizationPayloadKey).(*token.Payload)

	err = s.UserService.ChangeImage(c.Context(), file, payload.Email)
	if err != nil {
		return helper.Error(c, err)
	}

	return c.SendStatus(http.StatusOK)
}
