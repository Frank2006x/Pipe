package router

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/api/handler"
	"Frank2006x/Pipe/internal/api/middleware"
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/github"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RepositoryRouter(router *fiber.App, queries *sqlc.Queries, githubClient *github.Client, jwtMaker *auth.JwtMaker) {
	RepositoryService := service.NewRepositoryService(queries, githubClient)
	RepositoryHandler := handler.NewRepositoryHandler(RepositoryService)

	repositoryGroup := router.Group("/repositories")
	repositoryGroup.Use(middleware.AuthMiddleware(jwtMaker))
	repositoryGroup.Get("/hello", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	repositoryGroup.Post("/", RepositoryHandler.CreateRepository)
	repositoryGroup.Get("/", RepositoryHandler.ListAllRepositories)
}
