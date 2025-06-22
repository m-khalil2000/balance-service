# ---------- Stage 1: Build ----------
FROM golang:1.24 AS build

# Install git (required for go modules that use Git)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static build to avoid GLIBC dependency
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o balance-service .

# ---------- Stage 2: Distroless Minimal Runtime ----------
FROM gcr.io/distroless/static

COPY --from=build /app/balance-service /

ENTRYPOINT ["/balance-service"]
