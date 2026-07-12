package middleware

import (
	"Frank2006x/Pipe/internal/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func AuthMiddleware(jwt *auth.JwtMaker) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("jwt_token")

		log.Infof("Cookie Token: %s", token)
		if token == "" {
			return fiber.ErrUnauthorized
		}

		claims, err := jwt.ValidateJWT(token)
		log.Infof("Claims: %+v", claims)
		log.Infof("Err: %v", err)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals("user_id", claims.UserID)

		return c.Next()
	}
}
