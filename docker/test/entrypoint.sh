#!/bin/sh
set -e

echo "Waiting for PostgreSQL..."
until pg_isready -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME}; do
    echo "PostgreSQL is unavailable - sleeping"
    sleep 1
done

echo "PostgreSQL is up - executing tests"

./gophermarttest \
    -test.v \
    -test.run=^TestGophermart$ \
    -gophermart-binary-path=/app/gophermart \
    -gophermart-host=localhost \
    -gophermart-port=${MARKET_PORT} \
    -gophermart-database-uri="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
    -accrual-binary-path=/app/accrual \
    -accrual-host=localhost \
    -accrual-port=${ACCRUAL_PORT} \
    -accrual-database-uri="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
    2>&1 | tee /app/logs/gophermarttest.log 