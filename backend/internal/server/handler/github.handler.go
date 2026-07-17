package handler

import (
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

type GithubHandler struct {
	githubService *service.GithubService
}

func NewGithubHandler(githubService *service.GithubService) *GithubHandler {
	return &GithubHandler{
		githubService: githubService,
	}
}

func (h *GithubHandler) ListAllRepositories(c fiber.Ctx) error {

	userId := c.Locals("user_id").(int64)

	repositories, err := h.githubService.ListAllRepositories(c.Context(), userId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list repositories")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Repositories retrieved successfully",
		"repositories": repositories,
	})
}

func (h *GithubHandler) GetRepository(c fiber.Ctx) error {
	userId := c.Locals("user_id").(int64)
	owner := c.Params("owner")
	repo := c.Params("repo")

	repository, err := h.githubService.GetRepository(c.Context(), userId, owner, repo)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get repository")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Repository retrieved successfully",
		"repository": repository,
	})
}
