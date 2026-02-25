-- +goose Up
CREATE TABLE movies(
	id INTEGER PRIMARY KEY,
	created_at INT NOT NULL,
	updated_at INT NOT NULL,
	title TEXT NOT NULL,
	tmdb_id INT NOT Null,
	tmdb_url TEXT NOT NULL,
	poster_path TEXT NOT NULL,
	status TEXT NOT NULL
);

-- +goose Down
DROP TABLE movies;
