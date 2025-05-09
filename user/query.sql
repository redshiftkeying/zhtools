-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetPasswordByName :one
SELECT password FROM users
WHERE name = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: CreateUser :one
INSERT INTO users (
    name, password
) VALUES (
             ?, ?
         )
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
set name = ?,
    password = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;