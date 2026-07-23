package queue

import "context"

type Queue interface {
	PublishPipeline(ctx context.Context,pipelineId int64) error
	Close() error
}

type PipelineMessage struct{
	PipelineId int64 `json:"pipeline_id"`
	
}
	