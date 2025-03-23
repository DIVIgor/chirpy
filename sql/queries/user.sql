-- name: CreateUser :one
INSERT INTO users(id, email, hashed_password, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1;

-- name: ClearUsers :exec
DELETE FROM users;