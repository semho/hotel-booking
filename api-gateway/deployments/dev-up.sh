#!/bin/sh

echo "Starting API Gateway in development mode..."

# Загружаем переменные окружения
if [ -f "./deployments/.env" ]; then
    echo "Loading environment variables from .env file..."
    set -a
    . ./deployments/.env
    set +a
else
    echo "Warning: .env file not found in deployments directory"
fi

# Запускаем сервис
echo "Starting API Gateway service..."
APP_ENV=development go run cmd/main.go