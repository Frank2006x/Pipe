package router

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/api/handler"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RepositoryRouter(router *fiber.App, queries *sqlc.Queries) {
	RepositoryService := service.NewRepositoryService(queries)
	RepositoryHandler := handler.NewRepositoryHandler(RepositoryService)

	repositoryGroup := router.Group("/repositories")
	repositoryGroup.Post("/", RepositoryHandler.CreateRepository)

}
