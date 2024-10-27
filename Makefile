GOOSE_DRIVER='postgres'
GOOSE_DBSTRING='postgres://root@localhost:26257/layerg?sslmode=disable'
GOOSE_MIGRATION_DIR='./db/migrations'


all:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler --config .layerg-crawler.yaml

api:
	go build -ldflags -w
	chmod +x layerg-crawler
	./layerg-crawler api --config .layerg-crawler.yaml


migrate-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) up-to 20241024100253 # schema
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) up-to 20241027150911 # default data
migrate-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) down
migrate-reset:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) reset
