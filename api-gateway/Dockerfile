# Build stage
FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /api-gateway ./api-gateway/cmd/server

# Runtime stage
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /api-gateway /app/api-gateway

EXPOSE 8080

CMD ["/app/api-gateway"]