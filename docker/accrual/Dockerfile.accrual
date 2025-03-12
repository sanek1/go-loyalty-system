FROM golang:1.22-alpine AS builder

WORKDIR /app

# Установка необходимых зависимостей для сборки
RUN apk add --no-cache gcc musl-dev

# Сначала копируем только файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем все исходные файлы
COPY . .

# Проверка и сборка приложения
RUN ls -la ./cmd/accrual || echo "Директория cmd/accrual не существует"
RUN find . -name "accrual_linux_amd64" || echo "Файл accrual_linux_amd64 не найден"

# Если файл уже существует, используем его
RUN if [ -f ./cmd/accrual/accrual_linux_amd64 ]; then \
    cp ./cmd/accrual/accrual_linux_amd64 ./accrual; \
    else \
    # Иначе компилируем из исходников, если они есть
    if [ -d ./cmd/accrual ]; then \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./accrual ./cmd/accrual; \
    fi; \
    fi

# Проверка, что файл создан
RUN ls -la ./accrual || echo "Файл accrual не найден после сборки"

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

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/accrual ./accrual

# Установка прав на выполнение
RUN chmod +x ./accrual || echo "Файл ./accrual не найден"
RUN ls -la

EXPOSE ${ACCRUAL_PORT}

# Скрипт для запуска
COPY docker/accrual/entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"] 