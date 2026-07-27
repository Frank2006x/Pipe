package worker

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/queue"
	"Frank2006x/Pipe/internal/service"
	"context"
	"log"
	"time"
)

type PipelineWorker struct {
	pipelineService *service.PipelineService
	queue           queue.Queue
}

func NewPipelineWorker(queue queue.Queue, pipelineService *service.PipelineService) *PipelineWorker {
	return &PipelineWorker{
		pipelineService: pipelineService,
		queue:           queue,
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

	jobs, err := w.pipelineService.ListJobsByPipelineInternal(ctx, msg.PipelineId)
	if err != nil {
		return err
	}

	log.Printf("Pipeline %d has %d jobs", msg.PipelineId, len(jobs))

	for _, job := range jobs {
		log.Printf("Executing job %d", job.ID)
		err := w.ExecuteJob(ctx, job)
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

func (w *PipelineWorker) ExecuteJob(ctx context.Context, j sqlc.Job) error {
	now := time.Now()
	_, err := w.pipelineService.UpdateJobStatus(ctx, j.ID, sqlc.JobStatusRunning, &now, nil)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(2 * time.Second):
	}
	finished := time.Now()
	_, err = w.pipelineService.UpdateJobStatus(ctx, j.ID, sqlc.JobStatusSuccess, nil, &finished)
	if err != nil {
		return err
	}
	return nil
}
