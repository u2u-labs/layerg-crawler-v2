GOOSE_DRIVER='postgres'
GOOSE_DBSTRING='postgres://root@localhost:26257/layerg?sslmode=disable'
SYSTEM_MIGRATION_DIR='./db/migrations'
GENERATED_MIGRATION_DIR='./generated/migrations'
SERVICE_PORT='8084'

build:
	go build -ldflags -w
all:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler --config .layerg-crawler.yaml

query:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler query --config .layerg-crawler.yaml


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
	sqlc generate


# Generate Go helper files for all ABI files
gen-abi:
	@echo "Generating Go helper files from ABI..."
	@mkdir -p generated/abi_helpers
	@for file in abis/*.json; do \
		filename=$$(basename $$file .json); \
		echo "Processing $$filename..."; \
		go run main.go abigen -i $$file -o generated/abi_helpers/$${filename}_helpers.go; \
	done
	@echo "Done generating helpers!"

