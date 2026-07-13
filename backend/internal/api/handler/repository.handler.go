package handler

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/service"
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

func (h *RepositoryHandler) ImportRepository(c fiber.Ctx) error {
	var req struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
	}
	if err := c.Bind().Body(&req); err != nil {
		log.Errorf("Error parsing request body: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userId := c.Locals("userId").(int64)

	repository, err := h.repositoryService.ImportRepository(c.Context(), userId, req.Owner, req.Repo)

	if err != nil {
		log.Errorf("Error importing repository: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to import repository")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Repository imported successfully",
		"repository": repository,
	})
}

func (h *RepositoryHandler) ListAllRepositories(c fiber.Ctx) error {
	userId := c.Locals("userId").(int64)

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
	userId := c.Locals("userId").(int64)
	repoIdStr := c.Params("repoId")

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
