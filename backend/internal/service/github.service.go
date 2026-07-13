package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/github"
	"context"
)

type GithubService struct {
	githubClient *github.Client
	queries      *sqlc.Queries
}

func NewGithubService(githubClient *github.Client, queries *sqlc.Queries) *GithubService {
	return &GithubService{
		githubClient: githubClient,
		queries:      queries,
	}
}

func (s *GithubService) ListAllRepositories(ctx context.Context, userId int64) ([]github.Repository, error) {
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

func (s *GithubService) GetRepository(ctx context.Context, userId int64, owner string, repo string) (*github.Repository, error) {
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
