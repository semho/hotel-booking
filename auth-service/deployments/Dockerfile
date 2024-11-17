FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

RUN mkdir -p /app/bin

# Копируем модули из корня проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/auth-service ./auth-service/cmd/main.go

# Устанавливаем goose
RUN GOBIN=/app/bin go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

FROM alpine:latest
RUN apk add --no-cache postgresql16-client

WORKDIR /app

# Копируем из builder
COPY --from=builder /app/bin/auth-service .
COPY --from=builder /app/bin/goose /usr/local/bin/goose
COPY --from=builder /app/auth-service/deployments/migrations ./migrations
COPY --from=builder /app/auth-service/config/config.yaml ./config/
COPY --from=builder /app/auth-service/deployments/entrypoint.sh .

RUN chmod +x auth-service entrypoint.sh && \
    ls -la

RUN pwd && ls -la

EXPOSE 9092

CMD ["./entrypoint.sh"]