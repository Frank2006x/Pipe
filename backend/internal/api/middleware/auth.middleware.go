package middleware

import (
	"Frank2006x/Pipe/internal/auth"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(c fiber.Ctx, jwt auth.JwtMaker) error {

	token := c.Cookies("access_token")

	if token == "" {
		return fiber.ErrUnauthorized
	}

	claims, err := jwt.ValidateJWT(token)

	if err != nil {
		return fiber.ErrUnauthorized
	}

	c.Locals("user_id", claims.UserID)

	return c.Next()
}
