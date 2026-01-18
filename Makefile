include .env

export $(shell sed 's/=.*//' .env)

MIGRATIONS_DIR=./internal/migrations

create-migration:
	@if [ -z "$(name)" ]; then echo "Usage: make create-migration name=add_users"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" up

migrate-down-1:
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" down 1

migrate-down-all:

	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" down

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" version

migrate-up-n:
	@if [ -z "$(n)" ]; then echo "Usage: make migrate-up-n n=1"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" up $(n)

migrate-force:
	@if [ -z "$(v)" ]; then echo "Usage: make migrate-force v=3"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" force $(v)

