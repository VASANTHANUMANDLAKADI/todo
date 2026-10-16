# Stage 1: Download dependencies and prepare source code
FROM golang:1.26.3-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Stage 2: Run Go tests
FROM base AS test

RUN go test -count=1 ./...

# Stage 3: Build only after tests pass
FROM test AS builder

RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-api ./cmd/server

# Stage 4: Run the compiled application
FROM alpine:3.22 AS runtime

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /todo-api /app/todo-api

EXPOSE 8080

CMD ["/app/todo-api"]
