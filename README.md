# Bank Service Go

REST API банковского сервиса на Go.

Проект реализует регистрацию пользователей, JWT-аутентификацию, банковские счета, переводы, виртуальные карты, кредиты, график платежей, автоматическую обработку платежей, финансовую аналитику, интеграцию с ЦБ РФ и SMTP-уведомления.

## Стек

- Go
- gorilla/mux
- PostgreSQL
- lib/pq
- JWT
- bcrypt
- pgcrypto
- HMAC-SHA256
- logrus
- beevik/etree
- gomail.v2

## Основные возможности

- регистрация пользователей;
- аутентификация через JWT;
- создание банковских счетов;
- пополнение счетов;
- переводы между счетами;
- история операций;
- выпуск виртуальных карт;
- шифрование карточных данных через `pgcrypto`;
- HMAC-SHA256 для проверки целостности номера карты;
- bcrypt-хеширование CVV;
- оформление кредитов;
- генерация графика платежей;
- scheduler для автоматического списания кредитных платежей;
- штраф +10% за просрочку;
- аналитика доходов и расходов;
- расчет кредитной нагрузки;
- прогноз баланса;
- получение ключевой ставки ЦБ РФ через SOAP;
- SMTP-уведомления.

## Структура проекта

```text
cmd/bank-service       точка входа приложения
internal/config        конфигурация приложения
internal/db            подключение к PostgreSQL
internal/handlers      HTTP-обработчики
internal/middleware    JWT middleware
internal/models        модели и DTO
internal/repositories  работа с БД
internal/scheduler     scheduler кредитных платежей
internal/services      бизнес-логика
migrations             SQL-миграции
pkg/response           единый формат JSON-ответов
```

## Подготовка базы данных

Создайте базу данных:

```bash
psql -U postgres
```

```sql
CREATE DATABASE bank_service;
\q
```

Далее выполните SQL-файлы из папки `migrations` по порядку:

```text
001_create_users_table.sql
002_create_accounts_table.sql
003_create_transactions_table.sql
004_create_cards_table.sql
005_secure_cards_table.sql
006_create_credits_tables.sql
007_add_scheduler_indexes.sql
008_add_analytics_indexes.sql
```

## Конфигурация

Создайте локальный файл `.env` на основе `.env.example`.

Пример:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=bank_service
DB_SSLMODE=disable

JWT_SECRET=your_jwt_secret

CARD_PGP_KEY=your_card_pgp_key
CARD_HMAC_SECRET=your_card_hmac_secret

PAYMENT_SCHEDULER_INTERVAL_HOURS=12

BANK_RATE_MARGIN=5

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@example.com
SMTP_PASSWORD=your_smtp_password
SMTP_FROM=noreply@example.com
SMTP_ENABLED=false
```

Файл `.env` не должен попадать в Git.

## Запуск

```bash
go run ./cmd/bank-service
```

Проверка:

```bash
curl -s http://localhost:8080/health
```

Ожидаемый ответ:

```json
{"status":"ok"}
```

Проверка подключения к БД:

```bash
curl -s http://localhost:8080/health/db
```

Ожидаемый ответ:

```json
{"status":"ok","database":"available"}
```

## Регистрация

```bash
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","email":"user1@example.com","password":"user123"}'
```

Пример ответа:

```json
{
  "id": 1,
  "username": "user1",
  "email": "user1@example.com",
  "message": "user registered successfully"
}
```

## Аутентификация

```bash
curl -s -X POST http://localhost:8080/login \
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

Для дальнейших запросов можно сохранить токен:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user1@example.com","password":"user123"}' \
  | sed -E 's/.*"token":"([^"]+)".*/\1/')

echo "$TOKEN"
```

Проверка защищенного endpoint:

```bash
curl -s http://localhost:8080/me \
  -H "Authorization: Bearer $TOKEN"
```

## Счета

### Создать счет

```bash
curl -s -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"currency":"RUB"}'
```

### Получить свои счета

```bash
curl -s http://localhost:8080/accounts \
  -H "Authorization: Bearer $TOKEN"
```

### Пополнить счет

```bash
curl -s -X POST http://localhost:8080/accounts/1/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":5000}'
```

## Переводы

```bash
curl -s -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"fromAccountId":1,"toAccountId":2,"amount":1200}'
```

## История операций

```bash
curl -s http://localhost:8080/transactions \
  -H "Authorization: Bearer $TOKEN"
```

## Карты

### Выпуск карты

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

### Список карт

```bash
curl -s http://localhost:8080/cards \
  -H "Authorization: Bearer $TOKEN"
```

### Детали карты

```bash
curl -s http://localhost:8080/cards/1 \
  -H "Authorization: Bearer $TOKEN"
```

CVV повторно не возвращается.

## Кредиты

### Оформить кредит

```bash
curl -s -X POST http://localhost:8080/credits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"accountId":1,"principalAmount":100000,"interestRate":18,"termMonths":12}'
```

### Список кредитов

```bash
curl -s http://localhost:8080/credits \
  -H "Authorization: Bearer $TOKEN"
```

### График платежей

```bash
curl -s http://localhost:8080/credits/1/schedule \
  -H "Authorization: Bearer $TOKEN"
```

## Scheduler кредитных платежей

Scheduler запускается автоматически при старте приложения.

Интервал задается переменной:

```env
PAYMENT_SCHEDULER_INTERVAL_HOURS=12
```

Логика:

- ищет платежи с `payment_date <= CURRENT_DATE`;
- если денег хватает — списывает платеж;
- переводит платеж в статус `PAID`;
- уменьшает остаток долга;
- создает транзакцию `CREDIT_PAYMENT`;
- если денег не хватает — переводит платеж в статус `OVERDUE`;
- при первой просрочке увеличивает сумму платежа на 10%.

Для теста можно вручную сделать платеж текущим:

```bash
psql -U postgres -d bank_service
```

```sql
UPDATE payment_schedules
SET payment_date = CURRENT_DATE
WHERE id = 1;
```

После перезапуска приложения scheduler обработает платеж примерно через 5 секунд.

## Аналитика

### Доходы, расходы и кредитная нагрузка

```bash
curl -s http://localhost:8080/analytics \
  -H "Authorization: Bearer $TOKEN"
```

За конкретный месяц:

```bash
curl -s "http://localhost:8080/analytics?month=2026-05" \
  -H "Authorization: Bearer $TOKEN"
```

### Прогноз баланса

```bash
curl -s "http://localhost:8080/accounts/1/predict?days=30" \
  -H "Authorization: Bearer $TOKEN"
```

Максимальный период прогноза — 365 дней.

## Интеграция с ЦБ РФ

Endpoint получает ключевую ставку через SOAP API ЦБ РФ и добавляет банковскую маржу.

```bash
curl -s http://localhost:8080/rates/key \
  -H "Authorization: Bearer $TOKEN"
```

Пример ответа:

```json
{
  "centralBankRate": 16,
  "bankMargin": 5,
  "finalRate": 21,
  "message": "key rate received from Central Bank of Russia SOAP service"
}
```

## SMTP-уведомления

SMTP-настройки задаются через `.env`.

Если:

```env
SMTP_ENABLED=false
```

приложение не отправляет реальные письма, но endpoint работает безопасно для локальной проверки.

Тестовый запрос:

```bash
curl -s -X POST http://localhost:8080/notifications/test-email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"to":"user1@example.com"}'
```

Пример ответа:

```json
{
  "message": "test email processed successfully"
}
```

Если `SMTP_ENABLED=true`, scheduler также отправляет email-уведомления о кредитных платежах со статусами `PAID` и `OVERDUE`.

## Безопасность

В проекте реализовано:

- bcrypt-хеширование паролей;
- JWT-аутентификация;
- JWT middleware;
- добавление `userId` в context;
- проверка владельца счетов, карт и кредитов;
- bcrypt-хеширование CVV;
- шифрование номера карты через `pgcrypto`;
- шифрование срока действия карты через `pgcrypto`;
- HMAC-SHA256 для проверки целостности номера карты;
- параметризованные SQL-запросы;
- транзакции БД для переводов и кредитных операций.

## Финальная проверка проекта

```bash
go fmt ./...
go mod tidy
go test ./...
go run ./cmd/bank-service
```