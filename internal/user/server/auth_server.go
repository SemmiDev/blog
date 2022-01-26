package server

import (
	"github.com/SemmiDev/blog/internal/user/helper"
	"github.com/SemmiDev/blog/internal/user/service"

	"github.com/SemmiDev/blog/internal/common/memory"
	queryMemory "github.com/SemmiDev/blog/internal/user/query/memory"
	queryPostgresql "github.com/SemmiDev/blog/internal/user/query/postgresql"
	commandMemory "github.com/SemmiDev/blog/internal/user/repository/memory"
	commandPostgresql "github.com/SemmiDev/blog/internal/user/repository/postgresql"
	"github.com/SemmiDev/blog/internal/user/token"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

// AuthServer is the struct that contains UserService.
// it will be used to interact with the service.
type AuthServer struct {
	UserService service.UserService
}

// NewAuthServer creates a new AuthServer.
func NewAuthServer(db *pgxpool.Pool, tokenMaker token.Maker) (*AuthServer, error) {
	// for now, we're using memory repository for stores code verification.
	// in the future, we'll use redis, or other stores.
	m := memory.New()

	userServiceImpl := &service.UserServiceImpl{
		UserQuery:    queryPostgresql.NewUserQueryPostgresql(db),
		TokenQuery:   queryMemory.NewTokenQueryMemory(m),
		UserCommand:  commandPostgresql.NewUserCommandPostgresql(db),
		TokenCommand: commandMemory.NewTokenCommandMemory(m),
		TokenMaker:   tokenMaker,
	}

	return &AuthServer{UserService: userServiceImpl}, nil
}

// Mount mounts the auth server to the fiber app.
func (s *AuthServer) Mount(r fiber.Router) {
	r.Post("/register", s.RegisterHandler)
	r.Post("/authorize", s.AuthorizeHandler)
	r.Post("/password/reset", s.ResetPasswordHandler)
}

// RegisterHandler handles the registration of a new user.
// if the code in query param is empty,
// it will send a new verification code to the user's email.
// if the code in query param is not empty,
// it will verify the code has been sent to the user's email,
// and it will create a new user if user with email does not exist.
// if user with email already exists, it will return an error.
func (s *AuthServer) RegisterHandler(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	code := c.Query("code")
	if code != "" {
		userAuth, err := s.UserService.RegisterNewUser(c.Context(), code, name, password)
		if err != nil {
			return helper.Error(c, err)
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"data": userAuth,
		})
	}
	err := s.UserService.SendVerificationCode(c.Context(), email, "registration")
	if err != nil {
		return helper.Error(c, err)
	}
	return c.SendStatus(http.StatusOK)
}

// AuthorizeHandler handles user authorization/login
// based on email and password.
// it returns a token if user is authorized.
func (s *AuthServer) AuthorizeHandler(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	userAuth, err := s.UserService.Authorize(c.Context(), email, password)
	if err != nil {
		return helper.Error(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": userAuth,
	})
}

// ResetPasswordHandler handles password reset request.
// if code in query is empty, it will send verification
// code to the user's email for password reset.
// if code in query is not empty, it will reset user password.
func (s *AuthServer) ResetPasswordHandler(c *fiber.Ctx) error {
	email := c.FormValue("email")
	newPassword := c.FormValue("new_password")
	newConfirmPassword := c.FormValue("new_confirm_password")

	code := c.Query("code")
	if code != "" {
		err := s.UserService.ResetPassword(c.Context(), code, newPassword, newConfirmPassword)
		if err != nil {
			return helper.Error(c, err)
		}
		return c.SendStatus(http.StatusOK)
	}

	err := s.UserService.SendVerificationCode(c.Context(), email, "reset-password")
	if err != nil {
		return helper.Error(c, err)
	}
	return c.SendStatus(http.StatusOK)
}
