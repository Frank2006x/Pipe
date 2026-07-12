package router

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/api/handler"
	"Frank2006x/Pipe/internal/api/middleware"
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/config"
	"Frank2006x/Pipe/internal/github"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

func AuthRouter(router *fiber.App, queries *sqlc.Queries, githubClient *github.Client, config config.Config) {

	jwtMaker := auth.NewJwtMaker(config.JWT_SECRET)
	authService := service.NewAuthService(githubClient, jwtMaker, queries)
	authHandler := handler.NewAuthHandler(authService)
	authGroup := router.Group("/auth")

	authGroup.Get("/github", authHandler.GetRedirctLink)
	authGroup.Get("/github/callback", authHandler.Callback)
	authGroup.Get("/me", middleware.AuthMiddleware(jwtMaker), authHandler.GetUserInfo)
}
