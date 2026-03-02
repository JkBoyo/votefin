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

live/templ:
	templ generate --watch --cmd="go run ." --proxy="http://localhost:8080" --open-browser=false

live/tailwindcss:
	npx @tailwindcss/cli -i ./input.css -o ./assets/styles.css --watch

live:
	make -j2 live/templ live/tailwindcss
