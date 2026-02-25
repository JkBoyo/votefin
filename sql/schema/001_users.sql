-- +goose Up
CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	created_at INT NOT NULL,
	updated_at INT NOT NULL,
	jellyfin_user_id TEXT NOT NULL,
	username TEXT NOT NULL,
	is_admin INT NOT NULL
);

-- +goose Down
DROP TABLE users;
