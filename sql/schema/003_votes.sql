-- +goose Up
CREATE TABLE votes (
	id INT PRIMARY KEY,
	created_at INT NOT NULL,
	user_id INT NOT NULL REFERENCES users (id)
			ON DELETE CASCADE,
	movie_id INT NOT NULL REFERENCES users (id)
			ON DELETE CASCADE
);

-- +goose Down
DROP TABLE votes;
