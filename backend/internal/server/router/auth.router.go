package router

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/cache"
	"Frank2006x/Pipe/internal/github"
	"Frank2006x/Pipe/internal/server/handler"
	"Frank2006x/Pipe/internal/server/middleware"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

func AuthRouter(router *fiber.App, queries *sqlc.Queries, githubClient *github.Client, jwtMaker *auth.JwtMaker, cache *cache.Cache) {

	authService := service.NewAuthService(githubClient, jwtMaker, queries, cache)
	authHandler := handler.NewAuthHandler(authService)
	authGroup := router.Group("/auth")

	authGroup.Get("/github", authHandler.GetRedirctLink)
	authGroup.Get("/github/callback", authHandler.Callback)
	authGroup.Get("/me", middleware.AuthMiddleware(jwtMaker), authHandler.GetUserInfo)
	authGroup.Get("/logout", middleware.AuthMiddleware(jwtMaker), authHandler.Logout)
}
