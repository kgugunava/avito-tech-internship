# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o pr-reviewer-service ./cmd/main/main.go

# -------- Stage 2: Runtime --------
FROM alpine:3.18

# Устанавливаем необходимые зависимости (например для SSL)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/pr-reviewer-service .

# Порт, который слушает сервис
EXPOSE 8080

# Команда запуска
CMD ["./pr-reviewer-service"]
