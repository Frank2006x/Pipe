package worker

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/executor"
	"Frank2006x/Pipe/internal/queue"
	"Frank2006x/Pipe/internal/service"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PipelineWorker struct {
	pipelineService *service.PipelineService
	queue           queue.Queue
	executor        executor.Executor
}

func NewPipelineWorker(queue queue.Queue, pipelineService *service.PipelineService, exec executor.Executor) *PipelineWorker {
	return &PipelineWorker{
		pipelineService: pipelineService,
		queue:           queue,
		executor:        exec,
	}
}

func (w *PipelineWorker) Start(ctx context.Context) error {
	log.Println("[INFO] Pipeline Worker started")
	return w.queue.ConsumerPipeline(ctx, w.ProcessPipeline)
}

func (w *PipelineWorker) ProcessPipeline(ctx context.Context, msg queue.PipelineMessage) error {
	var err error
	defer func() {
		if err != nil {
			now := time.Now()
			dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer dbCancel()
			_, updateErr := w.pipelineService.UpdatePipelineStatus(dbCtx, msg.PipelineId, sqlc.PipelineStatusFailed, nil, &now)
			if updateErr != nil {
				log.Printf("Error while updating pipeline status: %v", updateErr)
			}
		}
	}()

	pipeline, err := w.pipelineService.GetPipelineInternal(ctx, msg.PipelineId)
	if err != nil {
		return err
	}

	if pipeline.Status != sqlc.PipelineStatusPending {
		return nil
	}

	now := time.Now()
	_, err = w.pipelineService.UpdatePipelineStatus(ctx, msg.PipelineId, sqlc.PipelineStatusRunning, &now, nil)
	if err != nil {
		return err
	}
	log.Printf("Pipeline %d status updated to running", msg.PipelineId)

	// Fetch repository information
	repo, err := w.pipelineService.GetRepositoryInternal(ctx, pipeline.RepositoryID)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Fetch github token if available
	token, _ := w.pipelineService.GetGithubTokenInternal(ctx, repo.UserID)

	// Create temporary workspace directory
	tmpDir, err := os.MkdirTemp("", "pipe-build-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() {
		log.Println("Cleaning workspace")
		_ = os.RemoveAll(tmpDir)
	}()

	// Clone repository
	log.Println("Cloning repository...")
	cloneUrl := repo.CloneUrl
	if token != "" && strings.HasPrefix(cloneUrl, "https://") {
		cloneUrl = strings.Replace(cloneUrl, "https://", "https://x-access-token:"+token+"@", 1)
	}

	cloneCmd := exec.CommandContext(ctx, "git", "clone", cloneUrl, tmpDir)
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository (%s): %s: %w", repo.FullName, string(output), err)
	}
	log.Println("Repository cloned")

	if pipeline.CommitSha != "" {
		checkoutCmd := exec.CommandContext(ctx, "git", "checkout", pipeline.CommitSha)
		checkoutCmd.Dir = tmpDir
		if output, err := checkoutCmd.CombinedOutput(); err != nil {
			log.Printf("Warning: failed to checkout commit %s: %s", pipeline.CommitSha, string(output))
		}
	}

	jobs, err := w.pipelineService.ListJobsByPipelineInternal(ctx, msg.PipelineId)
	if err != nil {
		return err
	}

	log.Printf("Pipeline %d has %d jobs", msg.PipelineId, len(jobs))

	for _, job := range jobs {
		log.Printf("Executing job %d (%s)", job.ID, job.Name)
		err := w.ExecuteJob(ctx, job, tmpDir)
		if err != nil {
			return err
		}
		log.Printf("Job %d executed successfully", job.ID)
	}

	finished := time.Now()
	_, err = w.pipelineService.UpdatePipelineStatus(ctx, msg.PipelineId, sqlc.PipelineStatusSuccess, nil, &finished)
	if err != nil {
		return err
	}

	return nil
}

func (w *PipelineWorker) ExecuteJob(ctx context.Context, j sqlc.Job, mountDir string) error {
	now := time.Now()
	_, err := w.pipelineService.UpdateJobStatus(ctx, j.ID, sqlc.JobStatusRunning, &now, nil)
	if err != nil {
		return err
	}

	execJob := executor.Job{
		Image:    "golang:1.24",
		MountDir: mountDir,
		WorkDir:  "/workspace",
		Commands: []string{"go mod download", "go build ./..."},
	}

	result, err := w.executor.Execute(ctx, execJob)
	finished := time.Now()

	if err != nil || result == nil || !result.Success {
		_, _ = w.pipelineService.UpdateJobStatus(ctx, j.ID, sqlc.JobStatusFailed, nil, &finished)
		if err != nil {
			return fmt.Errorf("job %d execution failed: %w", j.ID, err)
		}
		return fmt.Errorf("job %d execution failed with exit code %d", j.ID, result.ExitCode)
	}

	_, err = w.pipelineService.UpdateJobStatus(ctx, j.ID, sqlc.JobStatusSuccess, nil, &finished)
	if err != nil {
		return err
	}
	return nil
}
