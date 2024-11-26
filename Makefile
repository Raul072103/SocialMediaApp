include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations
DATABASE_PATH = "postgres://admin:adminpassword@localhost:5434/social?sslmode=disable"

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir ${MIGRATIONS_PATH} $(word 2, $(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_PATH) up $(word 2, $(MAKECMDGOALS))

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_PATH) down $(word 2, $(MAKECMDGOALS))
