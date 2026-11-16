# 💰 Cost Count In Time

Приложение, которое показывает, сколько времени нужно отработать, чтобы купить желаемую вещь.

## 🚀 Быстрый старт

### Backend (Go)

```bash
cd backend
go mod init costcountintimeapp
go run main.go
```

Сервер запустится на `http://localhost:8080`

### Frontend (React + Vite)

```bash
cd frontend
npm install
npm run dev
```

## 📦 Сборка для продакшена

```bash
# Frontend
cd frontend
npm run build

# Backend
cd backend
go build -o api
```

## 🌐 Деплой на Vercel

1. Установите Vercel CLI:
   ```bash
   npm i -g vercel
   ```

2. Задеплойте:
   ```bash
   vercel
   ```

## 🛠 Стек

| Компонент | Технология |
|-----------|------------|
| Backend | Go 1.21+ |
| Frontend | React 18 + TypeScript |
| Стили | TailwindCSS |
| Сборка | Vite |
| Деплой | Vercel |

## 📐 API

### POST /api/calculate

**Request:**
```json
{
  "salary": 50000,
  "salaryPeriod": "month",
  "price": 120000,
  "itemName": "iPhone 15"
}
```

**Response:**
```json
{
  "hours": 403.2,
  "days": 50.4,
  "formatted": "50 дн 3 ч",
  "itemName": "iPhone 15",
  "price": 120000
}
```

## 📝 Лицензия

MIT

## 🧱 Архитектура

Бэкенд разложен по слоям, зависимости идут только внутрь:

```
backend/
├── api.go                  HTTP-хендлер (точка входа, Vercel)
├── domain/                 сущности, value-объекты, интерфейсы репозиториев, доменные сервисы
├── application/            сервисы сценариев (calculate, karma) и DTO
└── infrastructure/         реализации хранилищ
```

Доменный слой не знает ни про HTTP, ни про базу: расчёт «сколько часов работы стоит вещь» и логика кармы живут в `domain/service`, а всё внешнее подключается через интерфейсы из `domain/repository`.

## 🧪 Тесты

```bash
cd backend
go test ./...
```

Тестами покрыты доменные сервисы и сервисы сценариев: `karma_calculator_test.go` (начисление кармы по правилам), `calculate_service_test.go` (расчёт часов и дней), `karma_service_test.go` (сценарии работы с событиями).

## 📝 Лицензия

MIT — см. [LICENSE](LICENSE).
