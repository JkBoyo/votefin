-- name: GetMovies :many
SELECT id, created_at, updated_at, title, tmdb_url, poster_path, status FROM movies;

-- name: InsertMovie :one
INSERT INTO movies (id, title, tmdb_url, poster_path, status )
VALUES (?, ?, ?, ?, ?)
RETURNING *;



