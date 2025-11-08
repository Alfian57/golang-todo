-include .env

DATABASE_URL ?= postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATION_DIR ?= migrations

kill-air:
	pkill -f "air"

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

migrate-up:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) up

migrate-drop:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) drop

swagger-generate:
	swag init -g ./main.go -o ./docs