-- name: GetMovies :many
SELECT id, title, url, poster_url, status, votes FROM movies;

-- name: InsertMovie :one
INSERT INTO movies (id, created_at, updated_at, title, url, poster_url, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);


