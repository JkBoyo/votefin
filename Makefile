include .env

dev-db-down:
	cd ./sql/schema/ &&  \
	goose sqlite3 ${DB_PATH} reset

dev-db-up:
	cd ./sql/schema/ && \
	goose sqlite3 ${DB_PATH} up

dev-db-up-down: dev-db-down dev-db-up


# dev-db-up-down:
# 	make dev-db-down
# 	make dev-db-up
