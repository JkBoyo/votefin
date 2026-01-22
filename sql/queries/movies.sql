-- name: GetMovies :many
SELECT id, created_at, updated_at, title, tmdb_id, tmdb_url, poster_path, status FROM movies;

-- name: GetMoviesSortedByVotes :many
SELECT m.id, m.created_at, m.updated_at, m.title, m.tmdb_id, m.tmdb_url, m.poster_path, m.status, COUNT(v.id) AS vote_count FROM movies m
INNER JOIN votes v on m.id = v.movie_id
GROUP BY m.id
ORDER BY vote_count DESC;

-- name: InsertMovie :one
INSERT INTO movies (created_at, updated_at, title, tmdb_id, tmdb_url, poster_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;



