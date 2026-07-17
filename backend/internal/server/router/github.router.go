package router

import (
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/server/handler"
	"Frank2006x/Pipe/internal/server/middleware"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

func GithubRouter(router *fiber.App, githubService *service.GithubService, jwtMaker *auth.JwtMaker) {
	githubHandler := handler.NewGithubHandler(githubService)

	githubGroup := router.Group("/github")
	githubGroup.Use(middleware.AuthMiddleware(jwtMaker))
	githubGroup.Get("/repositories", githubHandler.ListAllRepositories)
	githubGroup.Get("/repositories/:owner/:repo", githubHandler.GetRepository)
}
