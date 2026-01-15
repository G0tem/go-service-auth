.PHONY: up down rebuild swag test migrate-up migrate-down install-deps logs logs-app logs-db test-db-create test-db-drop test-db-recreate proto

# Установка зависимостей
install-deps:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Запуск приложения
run:
	docker compose -f docker-compose.yml up -d --build

# Остановка приложения
stop:
	docker compose -f docker-compose.yml down

# Пересборка и перезапуск
rerun:
	docker compose -f docker-compose.yml down
	docker compose -f docker-compose.yml build
	docker compose -f docker-compose.yml up -d

# Генерация Swagger документации
swag:
	~/go/bin/swag init -g main.go

# Генерация gRPC кода из proto файлов
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth.proto

# Запуск тестов
test: test-db-recreate
	go test -v -cover ./...

migrate-up:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/userdb?sslmode=disable" down

# Просмотр логов всех контейнеров
logs:
	docker compose -f docker-compose.yml logs -f

# Запуск тестов с отчетом о покрытии
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Создание тестовой БД
test-db-create:
	docker compose -f docker-compose.local.yml exec postgres psql -U postgres -c "CREATE DATABASE auth_test_db"

# Удаление тестовой БД
test-db-drop:
	docker compose -f docker-compose.local.yml exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS auth_test_db"

# Пересоздание тестовой БД
test-db-recreate: test-db-drop test-db-create 