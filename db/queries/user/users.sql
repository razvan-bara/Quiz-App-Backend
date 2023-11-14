-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY "firstName" AND "lastName";

-- name: CreateUser :one
INSERT INTO users (
    email, password, "firstName", "lastName", "isAdmin"
) VALUES (
             $1, $2, $3, $4, $5
         )
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;