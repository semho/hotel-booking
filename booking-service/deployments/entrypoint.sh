#!/bin/sh

echo "Starting entrypoint script..."

# Загружаем и экспортируем переменные из .env
set -a
source ./.env
set +a

# Проверяем подключение к БД
export PGPASSWORD="$DB_PASSWORD"
until psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
    echo "Postgres is unavailable - sleeping"
    sleep 1
done

echo "Postgres is up - executing migrations"

# Запускаем миграции
goose -dir ./migrations postgres "host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable" up

# Запускаем сервис
echo "Starting booking-service"
exec ./booking-service