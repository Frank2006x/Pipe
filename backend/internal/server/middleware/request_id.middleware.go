package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)
		log.Infof("Request id : %s", requestID)
		return c.Next()
	}
}
