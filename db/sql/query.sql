-- name: GetUserById :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash) VALUES (?, ? , ?) RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
