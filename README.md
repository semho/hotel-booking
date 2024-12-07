# Hotel Booking System

Микросервисное приложение для бронирования отелей, написанное на Go.

## Архитектура

Проект состоит из следующих микросервисов:

- **API Gateway**: Входная точка для всех клиентских запросов
- **Auth Service**: Управление аутентификацией и авторизацией
- **Hotel Service**: Управление отелями и номерами
- **Booking Service**: Управление бронированиями

## Технологии

- Go 1.23
- gRPC для межсервисного взаимодействия
- JWT для аутентификации
- PostgreSQL для хранения данных
- Docker и Docker Compose для контейнеризации
- Makefile для автоматизации задач

## Аутентификация

В проекте используется JWT-based аутентификация с парой токенов:

- **Access Token**: Короткоживущий токен для доступа к API
- **Refresh Token**: Долгоживущий токен для обновления access token

### Flow аутентификации:

1. **Логин/Регистрация**
   ```
   POST /api/v1/auth/login
   POST /api/v1/auth/register
   ```
   Возвращает:
   ```json
   {
     "access_token": "...",
     "refresh_token": "...",
     "user": {
       "id": "...",
       "email": "...",
       "firstName": "...",
       "lastName": "..."
     }
   }
   ```

2. **Использование API**
   - Добавьте access token в заголовок: `Authorization: Bearer <access_token>`

3. **Обновление токенов**
   ```
   POST /api/v1/auth/refresh
   Body: { "refresh_token": "..." }
   ```
   Возвращает новую пару токенов

## Запуск проекта

1. Клонируйте репозиторий:
   ```bash
   git clone <repository-url>
   ```

2. Установите зависимости:
   ```bash
   make install-deps
   ```

3. Создайте `.env` файлы для каждого сервиса на основе `.env.example`

4. Запустите сервисы:
   ```bash
   make run-services
   ```
   Или с пересборкой:
   ```bash
   make run-services-build
   ```
5. Примените миграции для каждого сервиса:
   ```bash
   make migrate-up SERVICE=auth-service DB_PORT=5431 DB_NAME=auth_service
   make migrate-up SERVICE=booking-service DB_PORT=5432 DB_NAME=booking_service
   make migrate-up SERVICE=room-service DB_PORT=5433 DB_NAME=room_service
   ```

Для остановки:
   ```bash
   make stop-services
   ```

## Разработка

### Доступные make команды:

- `make install-deps`: Установка зависимостей проекта
- `make run-services`: Запуск сервисов
- `make run-services-build`: Запуск сервисов с пересборкой
- `make stop-services`: Остановка сервисов
- `make generate`: Генерация proto файлов для всех сервисов
- `make migrate-up`: Применение миграций (требуется указать SERVICE)
- `make migrate-down`: Откат миграций (требуется указать SERVICE)
- `make migrate-status`: Статус миграций (требуется указать SERVICE)
- `make migrate-create`: Создание новой миграции (требуется указать SERVICE)

### Структура проекта:

```
.
├── api-gateway/         # API Gateway сервис
├── auth-service/        # Сервис аутентификации
├── booking-service/     # Сервис бронирований
├── hotel-service/       # Сервис отелей
├── pkg/                 # Общие пакеты
├── proto/              # Protobuf определения
├── docker-compose.yml  # Docker compose конфигурация
└── Makefile           # Команды для управления проектом
```

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация нового пользователя
- `POST /api/v1/auth/login` - Вход в систему
- `POST /api/v1/auth/refresh` - Обновление токенов

### Отели
- `GET /api/v1/hotels` - Список отелей
- `GET /api/v1/hotels/{id}` - Информация об отеле
- `POST /api/v1/hotels` - Создание отеля (Admin)
- `PUT /api/v1/hotels/{id}` - Обновление отеля (Admin)
- `DELETE /api/v1/hotels/{id}` - Удаление отеля (Admin)

### Бронирования
- `POST /api/v1/bookings` - Создание бронирования
- `GET /api/v1/bookings` - Список бронирований пользователя
- `GET /api/v1/bookings/{id}` - Информация о бронировании
- `DELETE /api/v1/bookings/{id}` - Отмена бронирования

## Лицензия

MIT