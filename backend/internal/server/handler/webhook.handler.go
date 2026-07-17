package handler

import (
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
}

func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (h *WebhookHandler) GitHubWebhook(c fiber.Ctx) error {
	event := c.Get("X-GitHub-Event")
	delivery := c.Get("X-GitHub-Delivery")

	var payload GitHubWebhookPayload
	if err := c.Bind().Body(&payload); err != nil {
		log.Errorf("Error parsing GitHub webhook payload: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid payload")
	}
	xHubSignature := c.Get("X-Hub-Signature-256")
	ok, err := h.webhookService.CheckSignature(c.Context(), 0, 0, xHubSignature, c.Body())
	if err != nil {
		log.Errorf("Error checking webhook signature: %v", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature")
	}
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature")
	}

	log.Infof(
		"GitHub webhook | event=%s delivery=%s repo=%s ref=%s commit=%s",
		event,
		delivery,
		payload.Repository.FullName,
		payload.Ref,
		payload.After,
	)
	return c.SendStatus(fiber.StatusOK)
}
