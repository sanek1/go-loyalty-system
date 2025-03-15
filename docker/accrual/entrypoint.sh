#!/bin/sh
set -e

echo "Проверка файловой системы..."
ls -la /app

echo "Waiting for PostgreSQL..."
until pg_isready -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME}; do
    echo "PostgreSQL is unavailable - sleeping"
    sleep 1
done

echo "PostgreSQL is up - executing command"

# Проверка наличия бинарного файла
if [ ! -f "./accrual" ]; then
    echo "ОШИБКА: Бинарный файл accrual не найден!"
    echo "Текущая директория:"
    pwd
    echo "Содержимое директории:"
    ls -la
    
    # Если есть файл accrual_linux_amd64, попробуем использовать его
    if [ -f "./cmd/accrual/accrual_darwin_arm64" ]; then
        echo "Найден accrual_linux_amd64, пробуем использовать его..."
        cp ./cmd/accrual/accrual_darwin_arm64 ./accrual
        chmod +x ./accrual
    else
        echo "Файл accrual_linux_amd64 также не найден!"
        exit 1
    fi
fi

# Запускаем accrual с правильными параметрами
echo "Выполняем ./accrual..."
exec ./accrual -a "${ACCRUAL_HOST}:${ACCRUAL_PORT}" \
    -d "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" 