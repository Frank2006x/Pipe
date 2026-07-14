package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/mapper"
	"context"
	"errors"
	"fmt"
)

var ErrRepoAlreadyExists = errors.New("repository has already been imported")
var ErrRepoNotFound = errors.New("repository not found")
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
		return sqlc.Repository{}, ErrRepoAlreadyExists
	}

	githubRepoParams := mapper.MapGitHubRepository(userId, githubRepo)

	newRepository, err := s.queries.CreateRepository(ctx, githubRepoParams)
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to create repository: %w", err)
	}

	return newRepository, nil
}

func (s *RepositoryService) ListAllRepositories(ctx context.Context, userId int64) ([]sqlc.Repository, error) {
	repositories, err := s.queries.ListRepositoriesByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}
	return repositories, nil
}

func (s *RepositoryService) GetRepository(ctx context.Context, userId int64, repoId int64) (sqlc.Repository, error) {
	repository, err := s.queries.GetRepositoryById(ctx, sqlc.GetRepositoryByIdParams{
		ID:     repoId,
		UserID: userId,
	})
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to get repository: %w", err)
	}
	return repository, nil
}

func (s *RepositoryService) DeleteRepository(ctx context.Context, userId int64, repoId int64) error {
	row,err := s.queries.DeleteRepository(ctx, sqlc.DeleteRepositoryParams{
		ID:     repoId,
		UserID: userId,
	})
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}
	if row == 0 {
		return ErrRepoNotFound
	}
	return nil
}
