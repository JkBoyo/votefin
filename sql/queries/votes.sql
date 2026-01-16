-- name: CreateVote :one
INSERT INTO votes (id, created_at, user_id, movie_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetMoviesByUserVotes :many
SELECT m.id, m.created_at, m.updated_at, m.title, m.tmdb_url, m.poster_path, m.status
FROM movies m
INNER JOIN votes v on m.id = v.movie_id
WHERE v.user_id = ?;

-- name: GetUsersByMoveisVoted :many
SELECT u.username
FROM users u
INNER JOIN votes v on u.id = v.user_id
WHERE v.movie_id = ?;

-- name: GetVotesCountPerUser :one
SELECT COUNT(id) FROM votes
WHERE user_id = ?;
