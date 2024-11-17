FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Копируем модули из корня проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./api-gateway/cmd/main.go

FROM alpine:latest

WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/api-gateway/config/config.yaml ./config/

EXPOSE 8080

CMD ["./main"]