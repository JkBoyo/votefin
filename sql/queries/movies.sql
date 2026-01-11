-- name: GetMovies :many
SELECT id, created_at, updated_at, title, tmdb_url, poster_path, status FROM movies;

-- name: GetMoviesSortedByVotes :many
SELECT m.id, m.created_at, m.updated_at, m.title, m.tmdb_url, m.poster_path, m.status, COUNT(v.id) AS vote_count FROM movies m
INNER JOIN votes v on m.id = v.movie_id
GROUP BY m.id
ORDER BY vote_count DESC;

-- name: InsertMovie :one
INSERT INTO movies (id, created_at, updated_at, title, tmdb_url, poster_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;



