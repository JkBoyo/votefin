-- name: AddUser :exec
INSERT INTO users (id, created_at, updated_at, jellyfin_user_id, username, is_admin)
VALUES (?,?,?,?,?,?);

-- name: GetUserByJellyID :one
SELECT * FROM users
WHERE jellyfin_user_id = ?;
