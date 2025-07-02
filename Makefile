.PHONY: migrate migrate-up migrate-down help

# Переменные
STORAGE_PATH=./storage/sso.db
MIGRATIONS_PATH=./cmd/migrations
MIGRATOR_PATH=./cmd/migrator

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

# Показать справку
help:
	@echo "Доступные команды:"
	@echo "  migrate        - Применить миграции к базе данных"
	@echo "  migrate-up     - Alias для migrate"
	@echo "  migrate-safe   - Создать директорию storage и применить миграции"
	@echo "  help           - Показать эту справку"
