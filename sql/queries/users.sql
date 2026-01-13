-- name: AddUser :one
INSERT INTO users (created_at, updated_at, jellyfin_user_id, username, is_admin)
VALUES (?,?,?,?,?)
RETURNING *;

-- name: GetUserByJellyID :one
SELECT * FROM users
WHERE jellyfin_user_id = ?;
