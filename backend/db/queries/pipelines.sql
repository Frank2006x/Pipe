-- name: CreatePipeline :one
INSERT INTO pipelines (
    repository_id,
    github_delivery_id,
    commit_sha,
    commit_message,
    branch,
    event_type,
    trigger_username
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7
)
RETURNING *;

-- name: GetPipelineById :one
SELECT * FROM pipelines
WHERE id = $1;

-- name: ListRepositoryPipelines :many
SELECT * FROM pipelines
WHERE repository_id = $1
ORDER BY created_at DESC;

-- name: UpdatePipelineStatus :one
UPDATE pipelines
SET
    status = $1,
    started_at = $2,    
    finished_at = $3,
    updated_at = NOW()
WHERE id = $4
RETURNING *;