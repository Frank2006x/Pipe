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
SELECT * FROM jobs
WHERE pipeline_id = $1;

-- name: UpdateJobStatus :one
UPDATE jobs
SET
    status = $1,
    started_at = $2,
    finished_at = $3,
    updated_at = NOW()
WHERE id = $4
RETURNING *;