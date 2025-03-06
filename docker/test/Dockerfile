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

# Копирование бинарных файлов
COPY dst1/gophermarttest ./gophermarttest
COPY cmd/accrual/accrual_linux_amd64 ./accrual
COPY cmd/gophermart/gophermart ./gophermart

# Установка прав на выполнение
RUN chmod +x ./gophermarttest ./accrual ./gophermart

# Создание директории для логов
RUN mkdir -p /app/logs

# Скрипт для запуска тестов
COPY docker/test/entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"] 