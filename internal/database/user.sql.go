// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const clearUsers = `-- name: ClearUsers :exec
DELETE FROM users
`

func (q *Queries) ClearUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, clearUsers)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users(id, email, hashed_password, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
RETURNING id, email, created_at, updated_at, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, created_at, updated_at, hashed_password, is_chirpy_red
FROM users
WHERE email = $1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET email = $2,
    hashed_password = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, created_at, updated_at, hashed_password, is_chirpy_red
`

type UpdateUserParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.ID, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgradeUserPlan = `-- name: UpgradeUserPlan :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING id, email, created_at, updated_at, hashed_password, is_chirpy_red
`

func (q *Queries) UpgradeUserPlan(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, upgradeUserPlan, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
