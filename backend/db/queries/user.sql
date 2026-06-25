-- name: CreateUser :one
INSERT INTO users (github_id, username,email, avatar_url)
VALUES ($1, $2, $3, $4) RETURNING *;


-- name: UpdateUser :one
UPDATE users SET username = $2, email = $3, avatar_url = $4, updated_at = NOW()
WHERE github_id = $1 RETURNING *;

