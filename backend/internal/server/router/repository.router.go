package router

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/server/handler"
	"Frank2006x/Pipe/internal/server/middleware"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RepositoryRouter(router *fiber.App, queries *sqlc.Queries, githubService *service.GithubService, jwtMaker *auth.JwtMaker, db *pgxpool.Pool) {
	RepositoryService := service.NewRepositoryService(queries, githubService, db)
	RepositoryHandler := handler.NewRepositoryHandler(RepositoryService)

	repositoryGroup := router.Group("/repositories")
	repositoryGroup.Use(middleware.AuthMiddleware(jwtMaker))
	repositoryGroup.Get("/hello", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	repositoryGroup.Post("/import", RepositoryHandler.ImportRepository)
	repositoryGroup.Get("/", RepositoryHandler.ListAllRepositories)
	repositoryGroup.Get("/:id", RepositoryHandler.GetRepositoryById)
	repositoryGroup.Delete("/:id", RepositoryHandler.DeleteRepository)
}
