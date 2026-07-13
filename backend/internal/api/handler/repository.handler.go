package handler

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

type RepositoryHandler struct {
	repositoryService *service.RepositoryService
}

func NewRepositoryHandler(repositoryService *service.RepositoryService) *RepositoryHandler {
	return &RepositoryHandler{
		repositoryService: repositoryService,
	}
}

func (h *RepositoryHandler) CreateRepository(c fiber.Ctx) error {
	var req sqlc.CreateRepositoryParams
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	repository, err := h.repositoryService.CreateRepository(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Repository created successfully",
		"repository": repository,
	})
}

