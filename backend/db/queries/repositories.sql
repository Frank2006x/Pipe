-- name: CreateRepository :one
INSERT INTO repositories 
    (user_id,
    github_repo_id,
    name,
    full_name,
    html_url,
    default_branch,
    private,
    webhook_id,
    is_active,
    owner,
    description,
    clone_url
    ,webhook_secret)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *;


-- name: GetRepositoryByGithubRepoID :one
SELECT * FROM repositories 
WHERE github_repo_id = $1;

-- name: GetRepositoryById :one
SELECT * FROM repositories 
WHERE id = $1 AND user_id = $2;

-- name: GetRepositoryByFullName :one
SELECT * FROM repositories 
WHERE full_name = $1;


-- name: ListRepositoriesByUser :many
SELECT * FROM repositories 
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateWebhookInfo :one
UPDATE repositories
SET 
    webhook_id = $1,
    webhook_secret = $2,
    is_active = $3,
    updated_at = NOW()
WHERE user_id = $4 AND github_repo_id = $5
RETURNING *;

-- name: UpdateRepositoryMetadata :one
UPDATE repositories
SET 
    name = $1,
    full_name = $2,
    owner = $3,
    description = $4,
    default_branch = $5,
    private = $6,
    html_url = $7,
    clone_url = $8,
    updated_at = NOW()
WHERE user_id = $9 AND github_repo_id = $10
RETURNING *;

-- name: SetRepositoryActive :one
UPDATE repositories
SET 
    is_active = $1,
    updated_at = NOW()
WHERE user_id = $2 AND github_repo_id = $3
RETURNING *;

-- name: GetRepositoryByWebhookId :one
SELECT * FROM repositories
WHERE webhook_id = $1;

-- name: DeleteRepository :execrows
DELETE FROM repositories WHERE id = $1 AND user_id = $2;


-- name: ExistsRepository :one
SELECT EXISTS (
    SELECT 1
    FROM repositories
    WHERE user_id = $1
      AND full_name = $2
);