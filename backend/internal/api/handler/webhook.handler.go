package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type WebhookHandler struct {
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
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
