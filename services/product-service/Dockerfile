# Build stage
FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /product-service ./services/product-service/cmd/server

# Runtime stage
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /product-service /app/product-service

EXPOSE 8082

CMD ["/app/product-service"]