# Этап сборки
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
    curl \
    postgresql-client \
    tzdata \
    musl-locales \
    musl-locales-lang
# Установка локалей
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

# Копирование бинарных файлов и миграций
COPY --from=builder /app/gophermart .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/cmd/accrual/accrual_linux_amd64 ./accrual
COPY --from=builder /app/dst1/gophermarttest.exe ./gophermarttest

# Установка прав на выполнение
RUN chmod +x ./gophermart ./accrual ./gophermarttest

# Создание директории для логов
RUN mkdir -p /app/logs

# Определение переменных окружения по умолчанию
ENV GOPHERMART_PORT=8080 \
    ACCRUAL_PORT=8081 \
    DB_URI="host=loyalty_db port=5432 user=admin password=admin dbname=admin sslmode=disable"

EXPOSE 8080 8081

# Скрипт для запуска
COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]