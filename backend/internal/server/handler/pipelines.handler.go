package handler

import (
	"Frank2006x/Pipe/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type PipelineHandler struct {
	pipelineService *service.PipelineService
}

func NewPipelineHandler(pipelineService *service.PipelineService) *PipelineHandler {
	return &PipelineHandler{
		pipelineService: pipelineService,
	}
}

func (h *PipelineHandler) GetRepositoryPipelines(c fiber.Ctx) error {
	repoIdStr := c.Params("id")
	repoId, err := strconv.ParseInt(repoIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid repository ID")
	}

	pipelines, err := h.pipelineService.ListRepositoryPipelines(c.Context(), repoId)
	if err != nil {
		log.Errorf("Error listing repository pipelines: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list repository pipelines")
	}

	return c.JSON(fiber.Map{
		"pipelines": pipelines,
	})
}

func (h *PipelineHandler) GetPipelineById(c fiber.Ctx) error {
	pipelineIdStr := c.Params("id")
	pipelineId, err := strconv.ParseInt(pipelineIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid pipeline ID")
	}

	pipeline, err := h.pipelineService.GetPipeline(c.Context(), pipelineId)
	if err != nil {
		log.Errorf("Error getting pipeline: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get pipeline")
	}

	return c.JSON(fiber.Map{
		"pipeline": pipeline,
	})
}

func (h *PipelineHandler) GetPipelineJobs(c fiber.Ctx) error {
	pipelineIdStr := c.Params("id")
	pipelineId, err := strconv.ParseInt(pipelineIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid pipeline ID")
	}

	jobs, err := h.pipelineService.ListJobsByPipeline(c.Context(), pipelineId)
	if err != nil {
		log.Errorf("Error listing pipeline jobs: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list pipeline jobs")
	}

	return c.JSON(fiber.Map{
		"jobs": jobs,
	})
}
