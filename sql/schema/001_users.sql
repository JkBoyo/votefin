-- +goose Up
CREATE TABLE users (
	id int NOT NULL,
	created_at text,
	updated_at text,
	name text,
	PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;
