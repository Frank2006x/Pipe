-- name: CreateToken :one
INSERT INTO github_tokens (user_id, access_token)
VALUES ($1, $2) RETURNING *;

-- name: UpdateToken :one
UPDATE github_tokens SET access_token = $2, updated_at = NOW()
WHERE user_id = $1 RETURNING *;