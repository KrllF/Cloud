include .env

APP_NAME = balancer
MAIN_FILE = cmd/app/main.go
LINT_FILE = .golangci.yaml

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN=$(PG_DSN)


build: lint
	go build -o $(APP_NAME) $(MAIN_FILE)

# Запуск приложения (со сборкой)
run: lint build
	./$(APP_NAME)

# Запуск приложения (без сборки)
run-only:
	./$(APP_NAME)

clear:
	rm $(APP_NAME)

# Обновляем зависимости
deps:
	go mod tidy

# линтер
lint:
	golangci-lint run -c $(LINT_FILE)

local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v
