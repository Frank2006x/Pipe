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
SELECT p.* FROM pipelines p
JOIN repositories r ON p.repository_id = r.id
WHERE p.id = $1 AND r.user_id = $2;

-- name: GetPipelineByIdInternal :one
SELECT * FROM pipelines
WHERE id = $1;

-- name: ListRepositoryPipelines :many
SELECT p.* FROM pipelines p
JOIN repositories r ON p.repository_id = r.id
WHERE p.repository_id = $1 AND r.user_id = $2
ORDER BY p.created_at DESC;

-- name: UpdatePipelineStatus :one
UPDATE pipelines
SET
    status = $1,
    started_at = COALESCE($2, started_at),
    finished_at = COALESCE($3, finished_at),
    updated_at = NOW()
WHERE id = $4
RETURNING *;