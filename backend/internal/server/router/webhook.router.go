package router

import (
	"Frank2006x/Pipe/internal/server/handler"

	"github.com/gofiber/fiber/v3"
)

func WebhookRouter(app *fiber.App) {
	webhookHandler := handler.NewWebhookHandler()
	app.Post("/webhooks/github", webhookHandler.GitHubWebhook);
}
