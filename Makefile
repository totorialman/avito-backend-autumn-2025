include .env

DB_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

migrate-status:
	goose -dir ./migrations postgres "$(DB_URL)" status

migrate-up:
	goose -dir ./migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_URL)" down

migrate-reset:
	goose -dir ./migrations postgres "$(DB_URL)" reset