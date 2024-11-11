include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations


.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) up

.PHONY: migrate-force
migrate-force:
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) force $(VERSION)


.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: drop-db
drop-db:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) drop

.PHONY: seed-db
seed-db:
	@go run cmd/migrate/seed/main.go