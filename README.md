# Bank Service Go.
REST API банковского сервиса на Go.

## Запуск
```bash
go run ./cmd/bank-service
```
### Проверка
```bash
curl http://localhost:8080/health
```

Ответ:
```json
{"status":"ok"}
```

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
```sql
CREATE DATABASE bank_service;
```

Проверка подключения к БД:
```bash
curl http://localhost:8080/health/db
```

Ожидаемый ответ:
```json
{"status":"ok","database":"available"}
```

---
## Регистрация пользователя

Endpoint:
```http
POST /register
```

Пример запроса:
```bash
curl -X POST http://localhost:8080/register \
 -H "Content-Type: application/json" \
 -d '{"username":"user1","email":"user1@example.com","password":"user123"}'
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

Пример запроса:
```bash
curl -X POST http://localhost:8080/login \
 -H "Content-Type: application/json" \
 -d '{"email":"user1@example.com","password":"user123"}'
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
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/login \
 -H "Content-Type: application/json" \
 -d '{"email":"user1@example.com","password":"user123"}' \
 | grep -o '"token":"[^"]*' \
 | cut -d'"' -f4)
curl http://localhost:8080/me \
 -H "Authorization: Bearer $TOKEN"
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

Endpoint:
```http
POST /accounts
```

Пример:
```bash
curl -X POST http://localhost:8080/accounts \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer $TOKEN" \
 -d '{"currency":"RUB"}'
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
Endpoint:
```http
GET /accounts
```

Пример:
```bash
curl http://localhost:8080/accounts \
 -H "Authorization: Bearer $TOKEN"
```
### Пополнение счета
Endpoint:
```http
POST /accounts/{accountId}/deposit
```

Пример:
```bash
curl -X POST http://localhost:8080/accounts/1/deposit \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer $TOKEN" \
 -d '{"amount":1500.50}'
```

Пример ответа:
```json
{
 "accountId": 1,
 "balance": 1500.5,
 "message": "account deposited successfully"
}
```

---
## Переводы и история операций

Все endpoints требуют JWT-токен:
```http
Authorization: Bearer <token>
```
### Перевод между счетами

Endpoint:
```http
POST /transfer
```

Пример:
```bash

curl -X POST http://localhost:8080/transfer \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer $TOKEN" \
 -d '{"fromAccountId":1,"toAccountId":2,"amount":1200}'
```

Пример ответа:
```json
{
 "transactionId": 1,
 "fromAccountId": 1,
 "toAccountId": 2,
 "amount": 1200,
 "fromBalance": 3800,
 "message": "transfer completed successfully"
}
```
### История операций

Endpoint:
```http
GET /transactions
```

Пример:
```bash
curl http://localhost:8080/transactions \
 -H "Authorization: Bearer $TOKEN"
```

Пример ответа:
```json
[
 {
  "id": 1,
  "userId": 1,
  "fromAccountId": 1,
  "toAccountId": 2,
  "transactionType": "TRANSFER",
  "amount": 1200,
  "createdAt": "2026-05-06T12:00:00Z"
 }
]
```

---
## Карты

Все endpoints требуют JWT-токен:
```http
Authorization: Bearer <token>
```

### Получение JWT-токена
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user1@example.com","password":"user123"}' \
  | sed -E 's/.*"token":"([^"]+)".*/\1/')

echo "$TOKEN"
```

### Выпуск виртуальной карты
```http
POST /cards
```

```bash
curl -s -X POST http://localhost:8080/cards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"accountId":1}'
```

Пример ответа:
```json
{
  "id": 1,
  "accountId": 1,
  "cardNumber": "2202123456789012",
  "maskedNumber": "220212******9012",
  "expiry": "05/2029",
  "cvv": "123",
  "message": "card created successfully. Save card number and CVV now; CVV will not be shown again."
}
```
Полный номер карты и CVV возвращаются только при выпуске карты.
### Получение списка карт
```http
GET /cards
```

```bash
curl -s http://localhost:8080/cards \
  -H "Authorization: Bearer $TOKEN"
```
В списке карт возвращается только маскированный номер.
### Получение деталей карты

```http
GET /cards/{cardId}
```

```bash
curl -s http://localhost:8080/cards/1 \
  -H "Authorization: Bearer $TOKEN"
```

Пример ответа:
```json
{
  "id": 1,
  "accountId": 1,
  "cardNumber": "2202123456789012",
  "maskedNumber": "220212******9012",
  "expiry": "05/2029",
  "status": "ACTIVE"
}
```
CVV не возвращается повторно.
### Безопасность карточных данных

В проекте используется следующая схема защиты:
- номер карты хранится в БД в зашифрованном виде через `pgcrypto`;
- срок действия карты хранится в БД в зашифрованном виде через `pgcrypto`;
- CVV хранится только как bcrypt-хеш;
- HMAC-SHA256 используется для проверки целостности номера карты;
- доступ к карте проверяется через JWT и `userId` владельца.

---

## Кредиты

Все endpoints требуют JWT-токен:
```http
Authorization: Bearer <token>
```

### Оформление кредита
```http
POST /credits
```

```bash
curl -s -X POST http://localhost:8080/credits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"accountId":1,"principalAmount":100000,"interestRate":18,"termMonths":12}'
```

Пример ответа:
```json
{
  "id": 1,
  "accountId": 1,
  "principalAmount": 100000,
  "interestRate": 18,
  "termMonths": 12,
  "monthlyPayment": 9168,
  "remainingAmount": 100000,
  "status": "ACTIVE",
  "message": "credit created successfully"
}
```

При оформлении кредита сумма кредита зачисляется на выбранный счет пользователя.

### Получение списка кредитов
```http
GET /credits
```

```bash
curl -s http://localhost:8080/credits \
  -H "Authorization: Bearer $TOKEN"
```

### Получение графика платежей
```http
GET /credits/{creditId}/schedule
```

```bash
curl -s http://localhost:8080/credits/1/schedule \
  -H "Authorization: Bearer $TOKEN"
```

Пример ответа:
```json
[
  {
    "id": 1,
    "creditId": 1,
    "paymentNumber": 1,
    "paymentDate": "2026-06-11T00:00:00Z",
    "amount": 9168,
    "principalPart": 7668,
    "interestPart": 1500,
    "status": "PLANNED",
    "createdAt": "2026-05-11T12:00:00Z"
  }
]
```

---
## Автоматическое списание кредитных платежей

В приложении реализован scheduler, который периодически обрабатывает платежи по кредитам.

Интервал задается переменной окружения:
```env
PAYMENT_SCHEDULER_INTERVAL_HOURS=12
```

Scheduler выполняет следующие действия:
- ищет платежи со сроком `payment_date <= CURRENT_DATE`;
- если на счете достаточно средств - списывает платеж;
- переводит платеж в статус `PAID`;
- уменьшает `remaining_amount` кредита;
- записывает операцию `CREDIT_PAYMENT` в историю транзакций;
- если средств недостаточно - переводит платеж в статус `OVERDUE`;
- при первой просрочке увеличивает сумму платежа на 10%.

### Проверка scheduler

Для теста можно вручную сделать ближайший платеж текущим:

```bash
psql -U postgres -d bank_service
```

```sql
UPDATE payment_schedules
SET payment_date = CURRENT_DATE
WHERE id = 1;
```

После перезапуска приложения scheduler обработает платеж примерно через 5 секунд.

Проверка графика:
```bash
curl -s http://localhost:8080/credits/1/schedule \
  -H "Authorization: Bearer $TOKEN"
```

Проверка истории операций:
```bash
curl -s http://localhost:8080/transactions \
  -H "Authorization: Bearer $TOKEN"
```