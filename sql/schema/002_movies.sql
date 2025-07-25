-- +goose Up
CREATE TABLE movies(
	id INT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	url TEXT NOT NULL,
	poster_url TEXT NOT NULL,
	status TEXT NOT NULL
);

-- +goose Down
DROP TABLE movies;
