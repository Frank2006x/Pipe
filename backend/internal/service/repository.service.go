package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/mapper"
	"Frank2006x/Pipe/internal/util/random"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRepoAlreadyExists = errors.New("repository has already been imported")
var ErrRepoNotFound = errors.New("repository not found")

type RepositoryService struct {
	queries       *sqlc.Queries
	githubService *GithubService
	db            *pgxpool.Pool
}

func NewRepositoryService(queries *sqlc.Queries, githubService *GithubService, db *pgxpool.Pool) *RepositoryService {
	return &RepositoryService{
		queries:       queries,
		githubService: githubService,
		db:            db,
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
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	exist, err := qtx.ExistsRepository(ctx, sqlc.ExistsRepositoryParams{
		UserID:   userId,
		FullName: githubRepo.FullName,
	})
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to check if repository exists: %w", err)
	}
	if exist {
		return sqlc.Repository{}, ErrRepoAlreadyExists
	}

	secret := random.GenerateRandomSecret()
	webhook, err := s.githubService.CreateWebhook(ctx, userId, owner, repo, secret)
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to create webhook: %w", err)
	}
	githubRepoParams := mapper.MapGitHubRepository(userId, githubRepo)

	_, err = qtx.CreateRepository(ctx, githubRepoParams)

	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to create repository: %w", err)
	}

	newRepository, err := qtx.UpdateWebhookInfo(ctx, sqlc.UpdateWebhookInfoParams{
		UserID:        userId,
		GithubRepoID:  githubRepo.ID,
		WebhookID:     pgtype.Int8{Int64: webhook.ID, Valid: true},
		WebhookSecret: pgtype.Text{String: secret, Valid: true},
		IsActive:      true,
	})

	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to update webhook info: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to commit transaction: %w", err)
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
	row, err := s.queries.DeleteRepository(ctx, sqlc.DeleteRepositoryParams{
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

func (s *RepositoryService) UpdateWebhookInfo(ctx context.Context, userId int64, githubRepoId int64, webhookId int64, webhookSecret string) (sqlc.Repository, error) {
	repository, err := s.queries.UpdateWebhookInfo(ctx, sqlc.UpdateWebhookInfoParams{
		UserID:        userId,
		GithubRepoID:  githubRepoId,
		WebhookID:     pgtype.Int8{Int64: webhookId, Valid: true},
		WebhookSecret: pgtype.Text{String: webhookSecret, Valid: true},
		IsActive:      true,
	})
	if err != nil {
		return sqlc.Repository{}, fmt.Errorf("failed to update webhook info: %w", err)
	}
	return repository, nil
}
