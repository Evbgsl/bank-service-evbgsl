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