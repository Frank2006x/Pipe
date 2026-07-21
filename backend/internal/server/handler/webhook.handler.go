package handler

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type WebhookHandler struct {
	webhookService  *service.WebhookService
	pipelineService *service.PipelineService
}

func NewWebhookHandler(webhookService *service.WebhookService, pipelineService *service.PipelineService) *WebhookHandler {
	return &WebhookHandler{
		webhookService:  webhookService,
		pipelineService: pipelineService,
	}
}

type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	HeadCommit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (h *WebhookHandler) GitHubWebhook(c fiber.Ctx) error {
	event := c.Get("X-GitHub-Event")
	if event == "ping" {
		log.Info("GitHub webhook ping received successfully")
		return c.SendStatus(fiber.StatusOK)
	}

	delivery := c.Get("X-GitHub-Delivery")
	xHubSignature := c.Get("X-Hub-Signature-256")

	var payload GitHubWebhookPayload
	if err := c.Bind().Body(&payload); err != nil {
		log.Errorf("Error parsing GitHub webhook payload: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid payload")
	}

	repo, err := h.webhookService.GetRepositoryByFullName(c.Context(), payload.Repository.FullName)
	if err != nil {
		log.Errorf("Error fetching repository: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch repository")
	}

	err = h.webhookService.CheckSignature(c.Context(), repo.UserID, repo.GithubRepoID, xHubSignature, c.Body())
	if err != nil {
		log.Errorf("Error checking webhook signature: %v", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature")
	}

	_, err = h.pipelineService.CreatePipeline(c.Context(), &service.CreatePipelineInput{
		RepositoryID:    repo.ID,
		DeliveryID:      delivery,
		CommitSHA:       payload.After,
		CommitMessage:   payload.HeadCommit.Message,
		Branch:          payload.Ref,
		EventType:       sqlc.GithubEvent(event),
		TriggerUsername: payload.Pusher.Name,
	})
	if err != nil {
		log.Errorf("Error creating pipeline: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create pipeline")
	}

	log.Infof(
		"GitHub webhook processed | event=%s delivery=%s repo=%s ref=%s commit=%s",
		event,
		delivery,
		payload.Repository.FullName,
		payload.Ref,
		payload.After,
	)

	return c.SendStatus(fiber.StatusOK)
}
