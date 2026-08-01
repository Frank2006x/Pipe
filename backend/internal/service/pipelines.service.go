package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/queue"
	"Frank2006x/Pipe/internal/util"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	queue   queue.Queue
}

func NewPipelineService(queries *sqlc.Queries, db *pgxpool.Pool, queue queue.Queue) *PipelineService {
	return &PipelineService{
		queries: queries,
		db:      db,
		queue:   queue,
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

	templates, err := qtx.ListRepositoryJobTemplatesInternal(ctx, input.RepositoryID)
	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		log.Printf("[INFO] No jobs configured for repository %d, skipping pipeline execution", input.RepositoryID)
		return nil, nil
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

	for _, tmpl := range templates {
		_, err := qtx.CreateJob(ctx, sqlc.CreateJobParams{
			PipelineID:       pipeline.ID,
			TemplateID:       pgtype.Int8{Int64: tmpl.ID, Valid: true},
			Status:           sqlc.JobStatusPending,
			Name:             tmpl.Name,
			OrderIndex:       tmpl.OrderIndex,
			Image:            tmpl.Image,
			WorkingDirectory: tmpl.WorkingDirectory,
			Commands:         tmpl.Commands,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	s.queue.PublishPipeline(ctx, pipeline.ID)
	return &pipeline, nil
}

type JobTemplateInput struct {
	Name             string   `json:"name"`
	OrderIndex       int32    `json:"order_index"`
	Image            string   `json:"image"`
	WorkingDirectory string   `json:"working_directory"`
	Commands         []string `json:"commands"`
}

func (s *PipelineService) ListRepositoryJobTemplates(ctx context.Context, userID, repositoryID int64) ([]sqlc.RepositoryJobTemplate, error) {
	return s.queries.ListRepositoryJobTemplates(ctx, sqlc.ListRepositoryJobTemplatesParams{
		RepositoryID: repositoryID,
		UserID:       userID,
	})
}

func (s *PipelineService) SaveRepositoryJobTemplates(ctx context.Context, userID, repositoryID int64, inputs []JobTemplateInput) ([]sqlc.RepositoryJobTemplate, error) {
	_, err := s.queries.GetRepositoryById(ctx, sqlc.GetRepositoryByIdParams{
		ID:     repositoryID,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("repository not found or access denied: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if err := qtx.DeleteRepositoryJobTemplatesByRepo(ctx, repositoryID); err != nil {
		return nil, err
	}

	var createdTemplates []sqlc.RepositoryJobTemplate
	for idx, input := range inputs {
		orderIdx := int32(idx)

		workDir := input.WorkingDirectory
		if workDir == "" {
			workDir = "/workspace"
		}

		tmpl, err := qtx.CreateRepositoryJobTemplate(ctx, sqlc.CreateRepositoryJobTemplateParams{
			RepositoryID:     repositoryID,
			Name:             input.Name,
			OrderIndex:       orderIdx,
			Image:            input.Image,
			WorkingDirectory: workDir,
			Commands:         input.Commands,
		})
		if err != nil {
			return nil, err
		}
		createdTemplates = append(createdTemplates, tmpl)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return createdTemplates, nil
}

func (s *PipelineService) GetPipeline(ctx context.Context, userId int64, id int64) (*sqlc.Pipeline, error) {
	pipeline, err := s.queries.GetPipelineById(ctx, sqlc.GetPipelineByIdParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

func (s *PipelineService) GetPipelineInternal(ctx context.Context, id int64) (*sqlc.Pipeline, error) {
	pipeline, err := s.queries.GetPipelineByIdInternal(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

func (s *PipelineService) GetRepositoryInternal(ctx context.Context, repositoryId int64) (*sqlc.Repository, error) {
	repo, err := s.queries.GetRepositoryByIdInternal(ctx, repositoryId)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (s *PipelineService) GetGithubTokenInternal(ctx context.Context, userId int64) (string, error) {
	token, err := s.queries.GetGithubToken(ctx, userId)
	if err != nil {
		return "", err
	}
	return token, nil
}


func (s *PipelineService) ListRepositoryPipelines(ctx context.Context, userId int64, repositoryId int64) ([]sqlc.Pipeline, error) {
	pipelines, err := s.queries.ListRepositoryPipelines(ctx, sqlc.ListRepositoryPipelinesParams{
		RepositoryID: repositoryId,
		UserID:       userId,
	})
	if err != nil {
		return nil, err
	}
	return pipelines, nil
}

func (s *PipelineService) ListJobsByPipeline(ctx context.Context, userId int64, pipelineId int64) ([]sqlc.Job, error) {
	jobs, err := s.queries.ListJobsByPipeline(ctx, sqlc.ListJobsByPipelineParams{
		PipelineID: pipelineId,
		UserID:     userId,
	})
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *PipelineService) ListJobsByPipelineInternal(ctx context.Context, pipelineId int64) ([]sqlc.Job, error) {
	jobs, err := s.queries.ListJobsByPipelineInternal(ctx, pipelineId)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *PipelineService) UpdatePipelineStatus(ctx context.Context, id int64, status sqlc.PipelineStatus, startedAt, finishedAt *time.Time) (*sqlc.Pipeline, error) {
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
	name             string
	status           sqlc.JobStatus
	orderIndex       int32
	image            string
	workingDirectory string
	commands         []string
}

func (s *PipelineService) createJob(ctx context.Context, queries *sqlc.Queries, pipelineID int64, input *createJobInput) (*sqlc.Job, error) {
	job, err := queries.CreateJob(ctx, sqlc.CreateJobParams{
		PipelineID:       pipelineID,
		Name:             input.name,
		Status:           input.status,
		OrderIndex:       input.orderIndex,
		Image:            input.image,
		WorkingDirectory: input.workingDirectory,
		Commands:         input.commands,
	})
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *PipelineService) createDefaultJobs(ctx context.Context, queries *sqlc.Queries, pipelineID int64) error {
	jobs := []createJobInput{
		{
			name:             "Build",
			status:           sqlc.JobStatusPending,
			orderIndex:       0,
			image:            "golang:1.24-alpine",
			workingDirectory: "/workspace",
			commands:         []string{"pwd", "ls -la", "go mod download", "go build ./..."},
		},
		{
			name:             "Test",
			status:           sqlc.JobStatusPending,
			orderIndex:       1,
			image:            "golang:1.24-alpine",
			workingDirectory: "/workspace",
			commands:         []string{"echo 'Running tests...'"},
		},
		{
			name:             "Deploy",
			status:           sqlc.JobStatusPending,
			orderIndex:       2,
			image:            "golang:1.24-alpine",
			workingDirectory: "/workspace",
			commands:         []string{"echo 'Deploying application...'"},
		},
	}

	for _, job := range jobs {
		_, err := s.createJob(ctx, queries, pipelineID, &job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PipelineService) UpdateJobStatus(ctx context.Context, id int64, status sqlc.JobStatus, startedAt, finishedAt *time.Time) (*sqlc.Job, error) {
	job, err := s.queries.UpdateJobStatus(ctx, sqlc.UpdateJobStatusParams{
		ID:         id,
		Status:     status,
		StartedAt:  util.TimestamptzOrNull(startedAt),
		FinishedAt: util.TimestamptzOrNull(finishedAt),
	})
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *PipelineService) UpdateJobResult(ctx context.Context, id int64, status sqlc.JobStatus, exitCode int32, logs string, finishedAt *time.Time) (*sqlc.Job, error) {
	job, err := s.queries.UpdateJobResult(ctx, sqlc.UpdateJobResultParams{
		ID:         id,
		Status:     status,
		ExitCode:   pgtype.Int4{Int32: exitCode, Valid: true},
		Logs:       logs,
		FinishedAt: util.TimestamptzOrNull(finishedAt),
	})
	if err != nil {
		return nil, err
	}
	return &job, nil
}
