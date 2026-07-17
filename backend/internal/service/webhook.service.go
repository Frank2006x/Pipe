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

func (s *WebhookService) CheckSignature(ctx context.Context, userId int64, githubRepoId int64, xHubSignature string, payload []byte) (bool, error) {

	secret, err := s.querier.GetRepoSecret(ctx, sqlc.GetRepoSecretParams{
		UserID:       userId,
		GithubRepoID: githubRepoId,
	})
	if err != nil {
		return false, fmt.Errorf("failed to fetch webhook secret: %w", err)
	}
	if xHubSignature == "" {
		return false, fmt.Errorf("X-Hub-Signature header is missing")
	}
	if !secret.Valid || secret.String == "" {
		return false, fmt.Errorf("webhook secret not found for userId: %d, githubRepoId: %d", userId, githubRepoId)
	}
	mac := hmac.New(sha256.New, []byte(secret.String))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	hmacCheck := hmac.Equal([]byte(expected), []byte(xHubSignature))
	if !hmacCheck {
		return false, fmt.Errorf("invalid signature: expected %s, got %s", expected, xHubSignature)
	}
	return true, nil
}
