-- +goose Up
CREATE TABLE movies(
	id INT PRIMARY KEY NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	title TEXT NOT NULL,
	tmdb_url TEXT NOT NULL,
	poster_path TEXT NOT NULL,
	status TEXT NOT NULL
);

-- +goose Down
DROP TABLE movies;
