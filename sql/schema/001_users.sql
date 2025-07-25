-- +goose Up
CREATE TABLE users (
	id INT PRIMARY KEY,
	created_at TEXT,
	updated_at TEXT,
	username TEXT
);

-- +goose Down
DROP TABLE users;
