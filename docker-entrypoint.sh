#!/bin/sh
set -e

# Определение переменных по умолчанию
DB_HOST=${DB_HOST:-loyalty_db}
DB_PORT=${DB_PORT:-5432}
DB_USER=${POSTGRES_USER:-admin}
DB_PASSWORD=${POSTGRES_PASSWORD:-admin}
DB_NAME=${POSTGRES_DB:-admin}

echo "Waiting for PostgreSQL..."
until pg_isready -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -f ./migrations/20250201101751_init_migrations.up.sql; do
    echo "PostgreSQL is unavailable - sleeping"
    sleep 1
done

echo "PostgreSQL is up - executing command"

# Запуск тестов если установлен флаг RUN_TESTS
if [ "$RUN_TESTS" = "true" ]; then
    echo "Running tests..."
    ./gophermarttest \
        -test.v \
        -test.run=^TestGophermart$ \
        -gophermart-binary-path=/app/gophermart \
        -gophermart-host=localhost \
        -gophermart-port=${GOPHERMART_PORT:-8080} \
        -gophermart-database-uri="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
        -accrual-binary-path=/app/accrual \
        -accrual-host=localhost \
        -accrual-port=${ACCRUAL_PORT:-8081} \
        -accrual-database-uri="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
        2>&1 | tee /app/logs/gophermarttest.log
fi

# Запуск основного приложения
exec ./gophermart \
    -a ":${GOPHERMART_PORT:-8080}" \
    -d "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
    -r "http://localhost:${ACCRUAL_PORT:-8081}"