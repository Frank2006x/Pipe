package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/util"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidEventType = errors.New("invalid github event type")
)

func isValidEventType(event sqlc.GithubEvent) bool {
	switch event {
	case sqlc.GithubEventPush,
		sqlc.GithubEventPullRequest,
		sqlc.GithubEventWorkflowDispatch,
		sqlc.GithubEventRelease,
		sqlc.GithubEventTag:
		return true
	default:
		return false
	}
}

type PipelineService struct {
	queries *sqlc.Queries
	db      *pgxpool.Pool
}

func NewPipelineService(queries *sqlc.Queries, db *pgxpool.Pool) *PipelineService {
	return &PipelineService{
		queries: queries,
		db:      db,
	}
}

type CreatePipelineInput struct {
	RepositoryID    int64
	DeliveryID      string
	CommitSHA       string
	CommitMessage   string
	Branch          string
	EventType       sqlc.GithubEvent
	TriggerUsername string
}

func (s *PipelineService) CreatePipeline(ctx context.Context, input *CreatePipelineInput) (*sqlc.Pipeline, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if !isValidEventType(input.EventType) {
		return nil, ErrInvalidEventType
	}

	pipeline, err := qtx.CreatePipeline(ctx, sqlc.CreatePipelineParams{
		RepositoryID:     input.RepositoryID,
		GithubDeliveryID: util.TextOrNull(input.DeliveryID),
		CommitSha:        input.CommitSHA,
		CommitMessage:    util.TextOrNull(input.CommitMessage),
		Branch:           input.Branch,
		EventType:        input.EventType,
		TriggerUsername:  util.TextOrNull(input.TriggerUsername),
	})
	if err != nil {
		return nil, err
	}

	err = s.createDefaultJobs(ctx, pipeline.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (s *PipelineService) GetPipeline(ctx context.Context, id int64) (*sqlc.Pipeline, error) {
	pipeline, err := s.queries.GetPipelineById(ctx, id)

	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

func (s *PipelineService) ListRepositoryPipelines(ctx context.Context, repositoryId int64) ([]sqlc.Pipeline, error) {
	pipelines, err := s.queries.ListRepositoryPipelines(ctx, repositoryId)
	if err != nil {
		return nil, err
	}
	return pipelines, nil
}

func (s *PipelineService) UpdatePipelineStatus(ctx context.Context, id int64, status sqlc.PipelineStatus, startedAt *time.Time, finishedAt *time.Time) (*sqlc.Pipeline, error) {
	pipeline, err := s.queries.UpdatePipelineStatus(ctx, sqlc.UpdatePipelineStatusParams{
		ID:         id,
		Status:     status,
		StartedAt:  util.TimestamptzOrNull(startedAt),
		FinishedAt: util.TimestamptzOrNull(finishedAt),
	})
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

type createJobInput struct {
	name       string
	status     sqlc.JobStatus
	orderIndex int32
}

func (s *PipelineService) createJob(ctx context.Context, pipelineID int64, input *createJobInput) (*sqlc.Job, error) {
	job, err := s.queries.CreateJob(ctx, sqlc.CreateJobParams{
		PipelineID: pipelineID,
		Name:       input.name,
		Status:     input.status,
		OrderIndex: input.orderIndex,
	})
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *PipelineService) createDefaultJobs(ctx context.Context, pipelineID int64) error {
	jobs := []createJobInput{
		{
			name:       "Build",
			status:     sqlc.JobStatusPending,
			orderIndex: 0,
		},
		{
			name:       "Test",
			status:     sqlc.JobStatusPending,
			orderIndex: 1,
		},
		{
			name:       "Deploy",
			status:     sqlc.JobStatusPending,
			orderIndex: 2,
		},
	}

	for _, job := range jobs {
		_, err := s.createJob(ctx, pipelineID, &job)
		if err != nil {
			return err
		}
	}

	return nil
}
