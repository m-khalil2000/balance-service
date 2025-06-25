# API Guide – Balance Service

This document describes the API exposed by the **Balance Service**. The service processes balance transactions (win/lose) for predefined users and allows retrieving their current balance.

---

## Base URL

```
http://localhost:8081
```

---

## Predefined Users

The following users are available by default (preloaded from the database):

| User ID | Initial Balance |
| ------- | --------------- |
| 1       | 100.00          |
| 2       | 200.00          |
| 3       | 300.00          |

---

## `GET /user/{userId}/balance`

### Description

Retrieve the current balance of a user.

---

### Successful Response

```json
{
  "userId": 1,
  "balance": "109.85"
}
```

---

### Error Responses

| Status Code | Cause          | Example Response              |
| ----------- | -------------- | ----------------------------- |
| `404`       | User not found | `{"error": "user not found"}` |

---

## `POST /user/{userId}/transaction`

### Description

Submit a transaction to increase or decrease the balance of a user.

Each `transactionId` must be unique. Repeating the same `transactionId` will result in a conflict error.

`transactionId` must be a valid UUID string (RFC 4122 format). This prevents duplicate transaction processing.

---

### Headers

| Header         | Value                       | Required |
| -------------- | --------------------------- | -------- |
| `Source-Type`  | `game`, `server`, `payment` | Yes      |
| `Content-Type` | `application/json`          | Yes      |

---

### Request Body

```json
{
  "state": "win",
  "amount": "10.15",
  "transactionId": "a6f4e8e2-3780-4c77-8f9b-234d53fd6eb6"
}
```

**Field Descriptions:**

- `state`: Either "win" or "lose"
- `amount`: String with up to 2 decimal places
- `transactionId`: Unique identifier per user

### Successful Response

```json
{
  "message": "transaction processed successfully",
  "oldBalance": "100.00",
  "newBalance": "110.15"
}
```

### Error Responses

| Status Code | Cause                                 | Example Response                             |
| ----------- | ------------------------------------- | -------------------------------------------- |
| `400`       | Invalid request body                  | `{"error": "invalid request payload"}`       |
| `404`       | User not found                        | `{"error": "user not found"}`                |
| `409`       | Duplicate transactionId               | `{"error": "transaction already processed"}` |
| `422`       | Insufficient balance or invalid state | `{"error": "insufficient balance"}`          |

### cURL Examples

**Valid win transaction**

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

**Valid lose transaction**

```bash
curl -X POST http://localhost:8081/user/2/transaction \
  -H "Source-Type: server" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "lose",
    "amount": "5.00",
    "transactionId": "a6f4e8e2-3780-4c77-8f9b-234d53fd6eb6"

  }'
```

**Duplicate transaction**

```bash
curl -X POST http://localhost:8081/user/2/transaction \
  -H "Source-Type: server" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "lose",
    "amount": "5.00",
    "transactionId": "a6f4e8e2-3780-4c77-8f9b-234d53fd6eb6"

  }'
```

**Invalid state**

```bash
curl -X POST http://localhost:8081/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "refund",
    "amount": "10.00",
    "transactionId": "a6f4e8e2-3780-4c77-8f9b-234d53fd6eb6"

  }'
```
