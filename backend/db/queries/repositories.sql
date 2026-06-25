-- name: CreateRepository :one
INSERT INTO repositories (name, github_url, default_branch)
VALUES ($1, $2, $3) RETURNING *;

-- name: FindByGithubUrl :one
SELECT * FROM repositories 
WHERE github_url = $1;

-- name: GetRepository :one
SELECT * FROM repositories 
WHERE name = $1;

-- name: ListRepositories :many
SELECT * FROM repositories 
ORDER BY created_at DESC;

-- name: DeleteRepository :exec
DELETE FROM repositories WHERE id = $1;