.PHONY: migrate migrate-up migrate-down run help

# Переменные
STORAGE_PATH=./storage/sso.db
MIGRATIONS_PATH=./cmd/migrations
MIGRATOR_PATH=./cmd/migrator
CONFIG_PATH=./config/local.yaml
SSO_PATH=./cmd/sso

# Применить миграции
migrate:
	go run $(MIGRATOR_PATH) --storage-path=$(STORAGE_PATH) --migrations-path=$(MIGRATIONS_PATH)

# Alias для migrate
migrate-up: migrate

# Создать директорию storage если она не существует
create-storage-dir:
	mkdir -p storage

# Запустить миграции с предварительным созданием директории
migrate-safe: create-storage-dir migrate

# Запустить сервер
run:
	go run $(SSO_PATH) --config=$(CONFIG_PATH)

# Показать справку
help:
	@echo "Доступные команды:"
	@echo "  migrate        - Применить миграции к базе данных"
	@echo "  migrate-up     - Alias для migrate"
	@echo "  migrate-safe   - Создать директорию storage и применить миграции"
	@echo "  run            - Запустить сервер"
	@echo "  help           - Показать эту справку"
