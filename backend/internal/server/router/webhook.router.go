package router

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/server/handler"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

func WebhookRouter(app *fiber.App, querier *sqlc.Queries, pipelineService *service.PipelineService) {
	webhookService := service.NewWebhookService(querier)
	webhookHandler := handler.NewWebhookHandler(webhookService, pipelineService)
	app.Post("/webhooks/github", webhookHandler.GitHubWebhook)
}
