# Go Auth Server

Сервер аутентификации на Go с использованием gRPC и SQLite.

## Описание

Это проект сервера аутентификации, который предоставляет следующие возможности:

- Регистрация пользователей
- Аутентификация с использованием JWT токенов
- Проверка прав администратора
- gRPC API для взаимодействия с клиентами
- SQLite база данных для хранения данных

## Структура проекта

- `cmd/` - точки входа приложения
  - `sso/` - основное приложение SSO сервера
  - `migrator/` - утилита для миграций базы данных
  - `migrations/` - SQL миграции
- `internal/` - внутренняя логика приложения
  - `app/` - инициализация приложения
  - `config/` - конфигурация
  - `domain/models/` - модели данных
  - `grpc/auth/` - gRPC сервер аутентификации
  - `lib/` - вспомогательные библиотеки
  - `services/auth/` - бизнес-логика аутентификации
  - `storage/` - слой работы с данными
- `config/` - файлы конфигурации
- `storage/` - файлы базы данных

## Требования

- Go 1.21+
- SQLite3

## Установка и запуск

1. Клонируйте репозиторий:

```bash
git clone https://github.com/HollyEllmo/go_auth_server.git
cd go_auth_server
```

2. Установите зависимости:

```bash
go mod download
```

3. Запустите миграции:

```bash
go run cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./cmd/migrations
```

4. Запустите сервер:

```bash
go run cmd/sso/main.go --config=./config/local.yaml
```

## Использование

Сервер предоставляет gRPC API для:

- Регистрации пользователей
- Входа в систему (получение JWT токена)
- Проверки прав администратора

## Конфигурация

Основная конфигурация находится в `config/local.yaml`.

## Миграции

Для управления схемой базы данных используется встроенный мигратор:

```bash
# Применить миграции
go run cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./cmd/migrations

# Применить миграции с определенной версии
go run cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./cmd/migrations --migrations-table=migrations
```

## Лицензия

MIT License
