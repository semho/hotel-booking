FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

RUN mkdir -p /app/bin

# Копируем модули из корня проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем room-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/room-service ./room-service/cmd/main.go

# Устанавливаем goose
RUN GOBIN=/app/bin go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest
RUN apk add --no-cache postgresql15-client

WORKDIR /app

# Копируем из builder
COPY --from=builder /app/bin/room-service .
COPY --from=builder /app/bin/goose /usr/local/bin/goose
COPY --from=builder /app/room-service/deployments/migrations ./migrations
COPY --from=builder /app/room-service/config/config.yaml ./config/
COPY --from=builder /app/room-service/deployments/entrypoint.sh .

RUN chmod +x room-service entrypoint.sh && \
    ls -la

# Проверяем наличие файлов
RUN pwd && ls -la

EXPOSE 8082 9092

CMD ["./entrypoint.sh"]