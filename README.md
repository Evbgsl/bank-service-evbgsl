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

