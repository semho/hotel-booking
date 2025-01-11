#!/bin/sh

echo "Starting development environment..."

# Загружаем переменные окружения
if [ -f "./deployments/.env" ]; then
    echo "Loading environment variables from .env file..."
    set -a
    . ./deployments/.env
    set +a
else
    echo "Warning: .env file not found in deployments directory"
fi

# Запускаем только БД
docker-compose -f deployments/docker-compose.dev.yml up -d

# Ждем готовности БД
echo "Waiting for PostgreSQL to be ready..."
while ! pg_isready -h localhost -p 5433 -U postgres -d room_service; do
    echo "PostgreSQL is unavailable - sleeping"
    sleep 1
done
echo "PostgreSQL is ready!"

# Применяем миграции, добавим проверку наличия goose
echo "Checking goose binary..."
ls -la ../bin/goose || echo "Goose binary not found in ../bin/"

# Применяем миграции
echo "Running migrations..."
../bin/goose -dir ./deployments/migrations postgres "host=localhost port=5433 user=postgres password=postgres dbname=room_service sslmode=disable" up

# Запускаем сервис
echo "Starting room service..."
APP_ENV=development go run cmd/main.go