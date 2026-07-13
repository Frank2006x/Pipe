package handler

import (
	"Frank2006x/Pipe/internal/service"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type RepositoryHandler struct {
	repositoryService *service.RepositoryService
}

func NewRepositoryHandler(repositoryService *service.RepositoryService) *RepositoryHandler {
	return &RepositoryHandler{
		repositoryService: repositoryService,
	}
}

func (h *RepositoryHandler) ImportRepository(c fiber.Ctx) error {
	var req struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
	}
	if err := c.Bind().Body(&req); err != nil {
		log.Errorf("Error parsing request body: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	log.Infof("Received import request: %+v", req)
	userId := c.Locals("user_id").(int64)
	log.Infof("Importing repository for user ID: %d, owner: %s, repo: %s\n", userId, req.Owner, req.Repo)
	repository, err := h.repositoryService.ImportRepository(c.Context(), userId, req.Owner, req.Repo)

	if err != nil {
		log.Errorf("Error importing repository: %v", err)
		if errors.Is(err, service.ErrRepoAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to import repository")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Repository imported successfully",
		"repository": repository,
	})
}

func (h *RepositoryHandler) ListAllRepositories(c fiber.Ctx) error {
	userId := c.Locals("user_id").(int64)

	repositories, err := h.repositoryService.ListAllRepositories(c.Context(), userId)
	if err != nil {
		log.Errorf("Error listing repositories: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list repositories")
	}

	return c.JSON(fiber.Map{
		"repositories": repositories,
	})
}

func (h *RepositoryHandler) GetRepositoryById(c fiber.Ctx) error {
	userId := c.Locals("user_id").(int64)
	repoIdStr := c.Params("id")

	repoId, err := strconv.ParseInt(repoIdStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid repository ID format",
		})
	}
	repository, err := h.repositoryService.GetRepository(c.Context(), userId, int64(repoId))
	if err != nil {
		log.Errorf("Error getting repository: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get repository")
	}

	return c.JSON(fiber.Map{
		"repository": repository,
	})
}
