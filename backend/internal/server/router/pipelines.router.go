package router

import (
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/server/handler"
	"Frank2006x/Pipe/internal/server/middleware"
	"Frank2006x/Pipe/internal/service"

	"github.com/gofiber/fiber/v3"
)

func PipelineRouter(router *fiber.App, pipelineService *service.PipelineService, jwtMaker *auth.JwtMaker) {
	pipelineHandler := handler.NewPipelineHandler(pipelineService)

	// Auth-protected endpoints
	pipelinesGroup := router.Group("/pipelines")
	pipelinesGroup.Use(middleware.AuthMiddleware(jwtMaker))
	pipelinesGroup.Get("/:id", pipelineHandler.GetPipelineById)
	pipelinesGroup.Get("/:id/jobs", pipelineHandler.GetPipelineJobs)

	// Repository pipelines & jobs endpoints
	repositoriesGroup := router.Group("/repositories")
	repositoriesGroup.Use(middleware.AuthMiddleware(jwtMaker))
	repositoriesGroup.Get("/:id/pipelines", pipelineHandler.GetRepositoryPipelines)
	repositoriesGroup.Get("/:id/jobs", pipelineHandler.GetRepositoryJobs)
	repositoriesGroup.Post("/:id/jobs", pipelineHandler.SaveRepositoryJobs)
}
