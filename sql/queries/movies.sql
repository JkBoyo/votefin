-- name: GetMovies :many
SELECT id, title, url, poster_url, status FROM movies;

-- name: InsertMovie :one
INSERT INTO movies (id, created_at, updated_at, title, description, url, poster_url, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;



