# Bank Service Go
REST API банковского сервиса на Go.

---
## Запуск
```bash
go run ./cmd/bank-service
```

## Проверка
```bash
curl http://localhost:8080/health
```
Ответ:  `{"status":"ok"}`

---

## Конфигурация

Приложение использует переменные окружения.

Пример файла конфигурации находится в `.env.example`.

Для локального запуска можно создать файл `.env`:
```env
APP_PORT=8080  
  
DB_HOST=localhost  
DB_PORT=5432  
DB_USER=postgres  
DB_PASSWORD=your_postgres_password  
DB_NAME=bank_service  
DB_SSLMODE=disable
```

Подготовка базы данных:
```SQL
CREATE DATABASE bank_service;
```

Проверка подключения к БД:
```bash
curl.exe http://localhost:8080/health/db
```

Ожидаемый ответ:
```JSON
{"status":"ok","database":"available"}
```

---

## Регистрация пользователя

Endpoint:
```http
POST /register
```

Пример запроса PowerShell:
```powershell
curl.exe -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{\"username\":\"user1\",\"email\":\"user1@example.com\",\"password\":\"user123\"}'
```

Пример успешного ответа:
```json
{
  "id": 1,
  "username": "user1",
  "email": "user1@example.com",
  "message": "user registered successfully"
}
```

Повторная регистрация с тем же `email` или `username` возвращает ошибку `409 Conflict`.

Пример ошибки:
```json
{
  "message": "user with this email or username already exists"
}
```

---
## Аутентификация

### Login

Endpoint:
```http
POST /login
```

Пример запроса PowerShell:
```powershell
curl.exe -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{\"email\":\"user1@example.com\",\"password\":\"user123\"}'
```

Пример ответа:
```json
{
  "token": "eyJ...",
  "tokenType": "Bearer",
  "expiresInSeconds": 86400
}
```

### Защищенный endpoint

Endpoint:
```http
GET /me
```

Запрос без токена вернет ошибку:
```json
{
  "message": "authorization header required"
}
```

Пример запроса с токеном:
```powershell
$response = curl.exe -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{\"email\":\"user1@example.com\",\"password\":\"user123\"}' | ConvertFrom-Json
$token = $response.token

curl.exe http://localhost:8080/me -H "Authorization: Bearer $token"
```

Пример ответа:
```json
{
  "userId": 1
}
```

---
## Счета

Все endpoints ниже требуют JWT-токен в заголовке:
```http
Authorization: Bearer <token>
```

### Создание счета
```http
POST /accounts
```

Пример PowerShell:
```powershell
curl.exe -X POST http://localhost:8080/accounts `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer $token" `
  -d '{\"currency\":\"RUB\"}'
```

Пример ответа:
```json
{
  "id": 1,
  "accountNumber": "40817810123456789012",
  "balance": 0,
  "currency": "RUB",
  "message": "account created successfully"
}
```

### Получение своих счетов
```http
GET /accounts
```

Пример:
```powershell
curl.exe http://localhost:8080/accounts -H "Authorization: Bearer $token"
```

### Пополнение счета
```http
POST /accounts/{accountId}/deposit
```

Пример:
```powershell
curl.exe -X POST http://localhost:8080/accounts/1/deposit `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer $token" `
  -d '{\"amount\":1500.50}'
```

Пример ответа:
```json
{
  "accountId": 1,
  "balance": 1500.5,
  "message": "account deposited successfully"
}
```

