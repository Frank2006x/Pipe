package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"context"
)

type RepositoryService struct {
	queries *sqlc.Queries
}

func NewRepositoryService(queries *sqlc.Queries) *RepositoryService {
	return &RepositoryService{
		queries: queries,
	}
}

func (s *RepositoryService) CreateRepository(ctx context.Context, req sqlc.CreateRepositoryParams) (sqlc.Repository, error) {

	repository, err := s.queries.CreateRepository(ctx, req)
	if err != nil {
		return sqlc.Repository{}, err
	}
	return repository, nil
}
