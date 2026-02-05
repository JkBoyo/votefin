-- name: CreateVote :one
INSERT INTO votes (created_at, user_id, movie_id, vote_count)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetMoviesByUserVotes :many
SELECT m.id, m.created_at, m.updated_at, m.title, m.tmdb_id, m.tmdb_url, m.poster_path, m.status, v.vote_count
FROM movies m
INNER JOIN votes v on m.id = v.movie_id
WHERE v.user_id = ?;

-- name: GetUsersByMoveisVoted :many
SELECT u.username
FROM users u
INNER JOIN votes v on u.id = v.user_id
WHERE v.movie_id = ?;

-- name: GetVotesCountPerUser :one
SELECT SUM(vote_count) FROM votes
WHERE user_id = ?;
