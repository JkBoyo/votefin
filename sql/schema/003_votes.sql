-- +goose Up
CREATE TABLE votes (
	id INTEGER PRIMARY KEY,
	created_at INT NOT NULL,
	updated_at INT NOT NULL,
	user_id INT NOT NULL REFERENCES users (id)
			ON DELETE CASCADE,
	movie_id INT NOT NULL REFERENCES movies (id)
			ON DELETE CASCADE,
	vote_count INT NOT NULL
);

-- +goose Down
DROP TABLE votes;
