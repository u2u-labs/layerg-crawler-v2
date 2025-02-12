GOOSE_DRIVER='postgres'
GOOSE_DBSTRING='postgres://root@localhost:26257/layerg?sslmode=disable'
SYSTEM_MIGRATION_DIR='./db/migrations'
GENERATED_MIGRATION_DIR='./generated/migrations'

build:
	go build -ldflags -w
all:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler --config .layerg-crawler.yaml

api:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler api --config .layerg-crawler.yaml

worker:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler worker --config .layerg-crawler.yaml

system-migrate-up:
	@echo "Running system migration..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(SYSTEM_MIGRATION_DIR) up
system-migrate-down:
	@echo "Reverting system migration..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(SYSTEM_MIGRATION_DIR) down

generated-migrate-up:
	@echo "Running generated migration..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GENERATED_MIGRATION_DIR) up
generated-migrate-down:
	@echo "Reverting generated migration..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GENERATED_MIGRATION_DIR) down
generated-migrate-reset:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GENERATED_MIGRATION_DIR) reset	

swag:
	swag init -g cmd/api_cmd.go -o ./docs

prepare:
	@echo "Running code generation flow..."
	go run cmd/prepare/main.go -schema=./schema.graphql -out=./generated -queries=./db
