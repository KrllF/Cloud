APP_NAME = balancer
MAIN_FILE = cmd/app/main.go
LINT_FILE = .golangci.yaml

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

