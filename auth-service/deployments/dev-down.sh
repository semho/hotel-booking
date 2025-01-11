#!/bin/sh

echo "Stopping development environment..."

# Останавливаем docker-compose
docker-compose -f deployments/docker-compose.dev.yml down

# Получаем порт из конфига
APP_ENV=development
PORT=$(go run cmd/tools/get_port.go)

if [ $? -eq 0 ] && [ ! -z "$PORT" ]; then
    # Находим и останавливаем только процесс Go сервиса
    pids=$(lsof -ti:$PORT)
    if [ ! -z "$pids" ]; then
        for pid in $pids; do
            # Проверяем, что это процесс Go
            if ps -p $pid -o comm= | grep -q "main"; then
                echo "Killing room service process on port $PORT (PID: $pid)..."
                kill $pid
            fi
        done
    else
        echo "Room service process not found on port $PORT"
    fi
else
    echo "Failed to get port from config"
fi