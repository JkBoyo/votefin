include .env

dev-db-reset:
	cd ./sql/schema/ &&  \
	goose sqlite3 ${DB_PATH} reset

dev-db-up:
	cd ./sql/schema/ && \
	goose sqlite3 ${DB_PATH} up

dev-db-down:
	cd ./sql/schema/ && \
	goose sqlite3 ${DB_PATH} down

dev-db-fr: dev-db-reset dev-db-up

dev-up:
	templ generate && \
	sqlc generate
	go run .
