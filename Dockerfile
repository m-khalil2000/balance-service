# ---------- Stage 1: Build ----------
FROM golang:1.24 AS build

# Install git (for go mod) and curl (if you want to test in builder stage too)
RUN apt-get update && apt-get install -y git curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static build to avoid runtime libc dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o balance-service .

# ---------- Stage 2: Lightweight Runtime ----------
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/balance-service /balance-service

# Run the service
ENTRYPOINT ["/balance-service"]
