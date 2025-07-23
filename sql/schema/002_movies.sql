-- +goose Up
CREATE TABLE movies(
	id INT,
	created_at TEXT,
	updated_at TEXT,
	votes BLOB,
);

-- +goose Down
DROP TABLE movies;
