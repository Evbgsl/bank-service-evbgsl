# Manual test scenario

## 1. Start server

```bash
go run ./cmd/bank-service
```

## 2. Health check

```bash
curl -s http://localhost:8080/health
curl -s http://localhost:8080/health/db
```

## 3. Register user

```bash
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","email":"user1@example.com","password":"user123"}'
```

## 4. Login

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user1@example.com","password":"user123"}' \
  | sed -E 's/.*"token":"([^"]+)".*/\1/')

echo "$TOKEN"
```

## 5. Create accounts

```bash
curl -s -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"currency":"RUB"}'

curl -s -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"currency":"RUB"}'
```

## 6. List accounts

```bash
curl -s http://localhost:8080/accounts \
  -H "Authorization: Bearer $TOKEN"
```

## 7. Deposit

```bash
curl -s -X POST http://localhost:8080/accounts/1/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":10000}'
```

## 8. Transfer

```bash
curl -s -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"fromAccountId":1,"toAccountId":2,"amount":1200}'
```

## 9. Transactions

```bash
curl -s http://localhost:8080/transactions \
  -H "Authorization: Bearer $TOKEN"
```

## 10. Create card

```bash
curl -s -X POST http://localhost:8080/cards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"accountId":1}'
```

## 11. List cards

```bash
curl -s http://localhost:8080/cards \
  -H "Authorization: Bearer $TOKEN"
```

## 12. Card payment

```bash
curl -s -X POST http://localhost:8080/cards/1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":750,"merchant":"Coffee Shop"}'
```

## 13. Create credit

```bash
curl -s -X POST http://localhost:8080/credits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"accountId":1,"principalAmount":100000,"interestRate":18,"termMonths":12}'
```

## 14. Credit schedule

```bash
curl -s http://localhost:8080/credits/1/schedule \
  -H "Authorization: Bearer $TOKEN"
```

## 15. Analytics

```bash
curl -s http://localhost:8080/analytics \
  -H "Authorization: Bearer $TOKEN"
```

## 16. Balance prediction

```bash
curl -s "http://localhost:8080/accounts/1/predict?days=30" \
  -H "Authorization: Bearer $TOKEN"
```

## 17. CBR key rate

```bash
curl -s http://localhost:8080/rates/key \
  -H "Authorization: Bearer $TOKEN"
```

## 18. Test email

```bash
curl -s -X POST http://localhost:8080/notifications/test-email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"to":"user1@example.com"}'
```