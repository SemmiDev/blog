package server

import (
	"github.com/SemmiDev/blog/internal/user/helper"
	"github.com/gofiber/fiber/v2"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware is a middleware that checks if the user is authenticated.
func (s *UserServer) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get the authorization header based on the authorizationHeaderKey
		authorizationHeader := c.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return helper.Error(c, helper.NewErr(helper.ErrAuthorizationHeaderKeyCode, authorizationHeaderKey))
		}

		// get the fields from the authorization header
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return helper.Error(c, helper.NewErr(helper.ErrAuthorizationHeaderFormatCode, authorizationHeaderKey))
		}

		// get the authorization type from the fields index 0
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return helper.Error(c, helper.NewErr(helper.ErrAuthorizationTypeBearerCode, authorizationTypeBearer))
		}

		// get the authorization payload from the fields index 1
		payload, err := s.TokenMaker.VerifyToken(fields[1])
		if err != nil {
			return helper.Error(c, helper.NewErr(helper.ErrAuthorizationInvalidTokenCode, "authorization_token"))
		}

		// set the payload to the context.
		// so that it can be used in the next middleware or handler.
		c.Context().SetUserValue(authorizationPayloadKey, payload)
		return c.Next()
	}
}
