include .env
export

.PHONY: migrate-up migrate-down run-api run-worker

# Накатить миграции в PostgreSQL
migrate-up:
	docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)" up

# Откатить миграции
migrate-down:
	docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)" down -all