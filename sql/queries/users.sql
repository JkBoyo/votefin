-- name: AddUser :exec
INSERT INTO users (id, created_at, updated_at, username)
VALUES (?,?,?,?);


