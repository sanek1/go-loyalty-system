# Loyalty System Service

Сервис для работы с системой лояльности, позволяющий накапливать и списывать баллы за заказы.

## Функциональность

- Регистрация и аутентификация пользователей
- Загрузка номеров заказов
- Начисление баллов за заказы
- Получение текущего баланса
- Списание баллов
- Получение информации о всех начислениях и списаниях
- Получение информации о всех загруженных заказах

## Технологии

- Go 1.22
- PostgreSQL
- Chi Router
- Zap Logger
- Docker & Docker Compose
- GitHub Actions CI/CD

## Требования

- Go 1.22 или выше
- PostgreSQL 15
- Docker & Docker Compose
- Make (опционально)

## Установка и запуск

### Локальный запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/your-username/loyalty-system.git
cd loyalty-system
```

2. Создайте файл конфигурации `.env`:
```bash
cat << EOF > .env
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=praktikum
ACCRUAL_SYSTEM_ADDRESS=http://localhost:8081
EOF
```

3. Запустите PostgreSQL:
```bash
docker-compose up -d postgres
```

4. Запустите приложение:
```bash
go run cmd/gophermart/main.go
```

### Запуск через Docker Compose

```bash
docker-compose up --build
```

## API Endpoints

### Аутентификация

- `POST /api/user/register` - Регистрация пользователя
- `POST /api/user/login` - Аутентификация пользователя

### Работа с заказами

- `POST /api/user/orders` - Загрузка номера заказа
- `GET /api/user/orders` - Получение списка заказов
- `GET /api/user/balance` - Получение текущего баланса
- `POST /api/user/balance/withdraw` - Списание баллов
- `GET /api/user/withdrawals` - Получение информации о выводе средств

## Структура проекта

```
.
├── cmd/
│   └── gophermart/
│       └── main.go
├── internal/
│   ├── app/
│   ├── entity/
│   ├── handler/
│   ├── usecase/
│   └── repo/
├── pkg/
│   ├── logger/
│   └── postgres/
├── migrations/
├── docker/
├── .github/
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Тестирование

### Запуск unit-тестов:
```bash
go test -v ./...
```

