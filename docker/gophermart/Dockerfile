FROM golang:1.22-alpine AS builder

WORKDIR /app

# Установка необходимых зависимостей
RUN apk add --no-cache gcc musl-dev

# Копирование и загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o gophermart ./cmd/gophermart

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Установка необходимых runtime зависимостей
RUN apk add --no-cache \
    postgresql-client \
    tzdata \
    musl-locales \
    musl-locales-lang

# Установка локалей
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

# Копирование бинарного файла и миграций
COPY --from=builder /app/gophermart .
COPY --from=builder /app/migrations ./migrations

# Установка прав на выполнение
RUN chmod +x ./gophermart

# Создание директории для логов
RUN mkdir -p /app/logs

EXPOSE ${MARKET_PORT}

# Скрипт для запуска
COPY docker/gophermart/entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"] 