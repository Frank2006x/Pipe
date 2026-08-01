-- name: CreateRepositoryJobTemplate :one
INSERT INTO repository_job_templates (
    repository_id,
    name,
    order_index,
    image,
    working_directory,
    commands
)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListRepositoryJobTemplates :many
SELECT t.* FROM repository_job_templates t
JOIN repositories r ON t.repository_id = r.id
WHERE t.repository_id = $1 AND r.user_id = $2
ORDER BY t.order_index ASC;

-- name: ListRepositoryJobTemplatesInternal :many
SELECT * FROM repository_job_templates
WHERE repository_id = $1
ORDER BY order_index ASC;

-- name: DeleteRepositoryJobTemplatesByRepo :exec
DELETE FROM repository_job_templates
WHERE repository_id = $1;
