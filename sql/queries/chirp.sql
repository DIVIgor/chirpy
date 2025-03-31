-- name: CreateChirp :one
INSERT INTO chirps(id, body, user_id, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
RETURNING *;

-- name: GetChirps :many
SELECT *
FROM chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT *
FROM chirps
WHERE ID = $1;

-- name: DeleteChirp :one
DELETE
FROM chirps
WHERE id = $1 AND user_id = $2
RETURNING *;