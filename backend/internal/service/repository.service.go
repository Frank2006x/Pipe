package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/mapper"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type RepositoryService struct {
	queries       *sqlc.Queries
	githubService *GithubService
}

func NewRepositoryService(queries *sqlc.Queries, githubService *GithubService) *RepositoryService {
	return &RepositoryService{
		queries:       queries,
		githubService: githubService,
	}
}

func (s *RepositoryService) CreateRepository(ctx context.Context, req sqlc.CreateRepositoryParams) (sqlc.Repository, error) {

	repository, err := s.queries.CreateRepository(ctx, req)
	if err != nil {
		return sqlc.Repository{}, err
	}
	return repository, nil
}

func (s *RepositoryService) ImportRepository(ctx context.Context, userId int64, owner string, repo string) (sqlc.Repository, error) {
	githubRepo, err := s.githubService.GetRepository(ctx, userId, owner, repo)
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to get repository from GitHub: %w", err)
	}

	exist, err := s.queries.ExistsRepository(ctx, sqlc.ExistsRepositoryParams{
		UserID:   userId,
		FullName: githubRepo.FullName,
	})
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to check if repository exists: %w", err)
	}
	if exist {
		return sqlc.Repository{}, fiber.NewError(fiber.StatusConflict, "Repository already exists")
	}

	githubRepoParams := mapper.MapGitHubRepository(userId, githubRepo)

	newRepository, err := s.queries.CreateRepository(ctx, githubRepoParams)
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to create repository: %w", err)
	}

	return newRepository, nil
}


