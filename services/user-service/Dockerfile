# Build stage
FROM golang:1.26.5 AS builder

WORKDIR /app

# Copy Go dependency files first for better layer caching.
COPY go.mod go.sum ./

RUN go mod download

# Copy the entire project.
COPY . .

# Build User Service.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /user-service ./services/user-service/cmd/server

# Runtime stage
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /user-service /app/user-service

EXPOSE 8081

CMD ["/app/user-service"]