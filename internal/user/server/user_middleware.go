package server

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (s *UserServer) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return Error(c, NewAuthorizationValidationError(ErrAuthorizationHeaderKey, authorizationHeaderKey))
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return Error(c, NewAuthorizationValidationError(ErrAuthorizationHeaderFormat, authorizationHeaderKey))
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return Error(c, NewAuthorizationValidationError(ErrAuthorizationTypeBearer, authorizationTypeBearer))
		}

		accessToken := fields[1]
		payload, err := s.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			return Error(c, NewAuthorizationValidationError(ErrAuthorizationInvalidToken, "authorization_token"))
		}

		c.Context().SetUserValue(authorizationPayloadKey, payload)
		return c.Next()
	}
}
