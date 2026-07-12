package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/github"
	"context"
)

type RepositoryService struct {
	queries      *sqlc.Queries
	githubClient *github.Client
}

func NewRepositoryService(queries *sqlc.Queries, githubClient *github.Client) *RepositoryService {
	return &RepositoryService{
		queries:      queries,
		githubClient: githubClient,
	}
}

func (s *RepositoryService) CreateRepository(ctx context.Context, req sqlc.CreateRepositoryParams) (sqlc.Repository, error) {

	repository, err := s.queries.CreateRepository(ctx, req)
	if err != nil {
		return sqlc.Repository{}, err
	}
	return repository, nil
}

func (s *RepositoryService) ListAllRepositories(ctx context.Context, userId int64) ([]github.Repository, error) {
	accessToken, err := s.queries.GetGithubToken(ctx, userId) // cache the access token in the future to avoid querying the database every time
	if err != nil {
		return nil, err
	}
	repositories, err := s.githubClient.ListRepositories(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return repositories, nil
}

func (s *RepositoryService) GetRepository(ctx context.Context, userId int64, owner string, repo string) (*github.Repository, error) {
	accessToken, err := s.queries.GetGithubToken(ctx, userId) // cache the access token in the future to avoid querying the database every time
	if err != nil {
		return nil, err
	}
	repository, err := s.githubClient.GetRepository(ctx, accessToken, owner, repo)
	if err != nil {
		return nil, err
	}
	return repository, nil

}
