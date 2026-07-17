package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type WebhookService struct {
	querier *sqlc.Queries
}

func NewWebhookService(querier *sqlc.Queries) *WebhookService {
	return &WebhookService{
		querier: querier,
	}
}

func VerifySignature(secret string, payload []byte, signature string) error {
	if signature == "" {
		return fmt.Errorf("signature is empty")
	}

	if secret == "" {
		return fmt.Errorf("webhook secret is empty")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)

	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func (s *WebhookService) CheckSignature(
	ctx context.Context,
	userID int64,
	githubRepoID int64,
	signature string,
	payload []byte,
) error {

	secret, err := s.querier.GetRepoSecret(ctx, sqlc.GetRepoSecretParams{
		UserID:       userID,
		GithubRepoID: githubRepoID,
	})
	if err != nil {
		return fmt.Errorf("get repository secret: %w", err)
	}

	if !secret.Valid || secret.String == "" {
		return fmt.Errorf("repository secret is not set: %w", err)
	}

	return VerifySignature(
		secret.String,
		payload,
		signature,
	)
}

func (s *WebhookService) GetUserAndRepoIdByFullName(
	ctx context.Context,
	fullName string,
) (int64, int64, error) {

	repo, err := s.querier.GetRepositoryByFullName(ctx, fullName)
	if err != nil {
		return 0, 0, fmt.Errorf("get repository by full name: %w", err)
	}

	return repo.UserID, repo.GithubRepoID, nil
}
