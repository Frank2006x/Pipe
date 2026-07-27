-- name: CreateJob :one
INSERT INTO jobs (
    pipeline_id,
    status,
    name,
    order_index
)
VALUES (
    $1,$2,$3,$4
)
RETURNING *;

-- name: ListJobsByPipeline :many
SELECT j.* FROM jobs j
JOIN pipelines p ON j.pipeline_id = p.id
JOIN repositories r ON p.repository_id = r.id
WHERE j.pipeline_id = $1 AND r.user_id = $2;

-- name: ListJobsByPipelineInternal :many
SELECT * FROM jobs
WHERE pipeline_id = $1;

-- name: UpdateJobStatus :one
UPDATE jobs
SET
    status = $1,
    started_at = COALESCE($2, started_at),
    finished_at = COALESCE($3, finished_at),
    updated_at = NOW()
WHERE id = $4
RETURNING *;