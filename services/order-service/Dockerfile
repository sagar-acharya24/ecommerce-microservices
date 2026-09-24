# Build stage
FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /order-service ./services/order-service/cmd/server

# Runtime stage
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /order-service /app/order-service

EXPOSE 8083

CMD ["/app/order-service"]