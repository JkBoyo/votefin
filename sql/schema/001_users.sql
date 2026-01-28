-- +goose Up
CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	jellyfin_user_id TEXT NOT NULL,
	username TEXT NOT NULL,
	is_admin INT NOT NULL
);

-- +goose Down
DROP TABLE users;
