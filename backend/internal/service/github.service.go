package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/config"
	"Frank2006x/Pipe/internal/github"
	"context"
	"fmt"
)

type GithubService struct {
	githubClient *github.Client
	queries      *sqlc.Queries
	config       *config.Config
}

func NewGithubService(githubClient *github.Client, queries *sqlc.Queries, config *config.Config) *GithubService {
	return &GithubService{
		githubClient: githubClient,
		queries:      queries,
		config:       config,
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

func (s *GithubService) CreateWebhook(ctx context.Context, userId int64, owner string, repo string, secret string) (*github.Webhook, error) {
	accessToken, err := s.queries.GetGithubToken(ctx, userId)

	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub token: %w", err)
	}
	webhookConfig := github.WebhookConfig{
		URL:         s.config.WEBHOOK_BASE_URL,
		ContentType: "json",
		Secret:      secret,
		InsecureSSL: "0",
	}
	webHook, err := s.githubClient.CreateWebhook(ctx, accessToken, owner, repo, webhookConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}
	return webHook, nil
}
