# Balance Service

The **Balance Service** is a REST API built with Go and PostgreSQL that tracks user balances and processes idempotent transactions. It supports basic `win`/`lose` transactions and ensures safe concurrency and validation.

---

## Features

- Get current balance of predefined users
- Process `win` or `lose` transactions
- Idempotent: no duplicate transactions with same `transactionId`
- Input validation and error handling
- Functional test suite with logging and concurrency checks
- Dockerized setup with PostgreSQL

---

## Tech Stack

- Go 1.21+
- PostgreSQL 15
- Gin (HTTP framework)
- `shopspring/decimal` for precise money math
- Docker & docker-compose

---

## Predefined Users

| User ID | Initial Balance |
| ------- | --------------- |
| 1       | 10000.00        |
| 2       | 20000.00        |
| 3       | 30000.00        |

---

## Running the Project

### Clone the repository

```bash
git clone https://github.com/your-username/balance-service.git
cd balance-service
```

## Start the Service

```bash
docker-compose up -d
```

Service will be available at:

```bash
http://localhost:8081
```

---

## API Endpoints

### `GET /user/{userId}/balance`

Fetch the current balance of a user.

#### Response

```json
{
  "userId": 1,
  "balance": "109.85"
}
```

---

### `POST /user/{userId}/transaction`

Submit a transaction (`win` or `lose`).

#### Required Headers

| Header       | Example                 |
| ------------ | ----------------------- |
| Source-Type  | `game, server, payment` |
| Content-Type | `application/json`      |

#### Request Body

```json
{
  "state": "win",
  "amount": "10.00",
  "transactionId": "uuid-string"
}
```

#### Response

```json
{
  "message": "transaction processed successfully",
  "oldBalance": "100.00",
  "newBalance": "110.00"
}
```

---

## Running Tests

Run the functional test suite:

```bash
go run test/functional.go
```

- A summary will be printed in the terminal
- Full logs saved to `test_suite.log`

---

## Project Structure

```bash
.
├── internal/
│   ├── handlers/       # HTTP route handlers
│   └── storage/        # PostgreSQL logic
├── pkg/
│   └── models/         # Payloads & response types
├── test/
│   └── functional.go   # Functional tests
├── go.mod / go.sum
├── docker-compose.yml
└── main.go
```

---

## Sample cURL

```bash
curl -X POST http://localhost:8081/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "10.00",
    "transactionId": "a6f4e8e2-3780-4c77-8f9b-234d53fd6eb6"
  }'
```
