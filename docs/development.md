# E-Commerce Microservices — Development Guide

## 1. Overview

This document explains how to set up, run, test, and develop the
E-Commerce Microservices project.

The project is implemented using Golang and follows a microservices
architecture.

Phase 1 uses:

- Golang
- REST APIs
- PostgreSQL
- Redis
- Docker
- Docker Compose
- JWT authentication
- Gin
- GORM

---

## 2. Prerequisites

The development environment requires:

- Git
- Go 1.26+
- Docker
- Docker Compose

Verify the installations:

```bash
go version
git --version
docker --version
docker compose version
3. Clone the Repository

Clone the repository:

git clone https://github.com/sagar-acharya24/ecommerce-microservices.git

Enter the project directory:

cd ecommerce-microservices
4. Environment Configuration

Create the local environment file:

cp .env.example .env

The .env file contains local development configuration.

The .env file should not be committed to Git.

The repository contains .env.example as a template.

5. Local Development

The services can be run individually during development.

User Service
go run ./services/user-service/cmd/server

Default port:

8081
Product Service
go run ./services/product-service/cmd/server

Default port:

8082
Order Service
go run ./services/order-service/cmd/server

Default port:

8083
API Gateway
go run ./api-gateway/cmd/server

Default port:

8080

The API Gateway is the recommended entry point for API requests.

6. Running with Docker Compose

Docker Compose runs the complete Phase 1 stack.

Start the stack:

docker compose -f deployments/docker-compose.yml up -d --build

Check the running containers:

docker compose -f deployments/docker-compose.yml ps

Expected services:

postgres
redis
user-service
product-service
order-service
api-gateway
7. Check API Gateway Health

The API Gateway exposes the public health endpoint:

curl http://localhost:8080/health

Expected response:

{
  "status": "ok"
}
8. View Logs

View logs for the complete stack:

docker compose -f deployments/docker-compose.yml logs

Follow logs:

docker compose -f deployments/docker-compose.yml logs -f

View logs for a specific service:

docker compose -f deployments/docker-compose.yml logs -f user-service
docker compose -f deployments/docker-compose.yml logs -f product-service
docker compose -f deployments/docker-compose.yml logs -f order-service
docker compose -f deployments/docker-compose.yml logs -f api-gateway
9. Stop Docker Compose

Stop the running containers:

docker compose -f deployments/docker-compose.yml down

This stops and removes the containers but keeps the named volumes.

Database and Redis data are persisted through Docker volumes.

Avoid using:

docker compose -f deployments/docker-compose.yml down -v

unless the intention is to remove the persistent database and Redis volumes.

10. Build the Go Project

Download dependencies:

go mod download

Tidy dependencies:

go mod tidy

Build all packages:

go build ./...
11. Run Unit Tests

Run the complete test suite:

go test ./...

Run tests with verbose output:

go test ./... -v

Run tests with race detection:

go test -race ./...
12. Run Static Analysis

Run Go vet:

go vet ./...

The project should pass go vet before changes are committed.

13. Product Service Tests

Run Product Service tests:

go test ./services/product-service/... -v

Product Service tests cover:

Product repository
Product service
Product handlers
Redis cache
14. Order Service Tests

Run Order Service tests:

go test ./services/order-service/... -v

Order Service tests cover:

Product client
Order repository
Order service
Order handlers
Authentication and authorization behavior
15. User Service Tests

Run User Service tests:

go test ./services/user-service/... -v

User Service tests cover:

User registration
Login
Authentication
User operations
Password handling
16. API Gateway

Run API Gateway:

go run ./api-gateway/cmd/server

The Gateway provides the external API entry point.

Default address:

http://localhost:8080

Backend services should not be used as the public API entry point when running
the Docker deployment.

17. Authentication Flow

User authentication follows this flow:

Client
  |
  | Register / Login
  v
API Gateway
  |
  v
User Service
  |
  | JWT
  v
Client

For protected requests:

Client
  |
  | Authorization: Bearer <JWT>
  v
API Gateway
  |
  | Validate JWT
  | Extract user_id
  | Set trusted X-User-ID
  v
Backend Service

The Gateway does not trust a client-provided X-User-ID.

18. Product Cache

Product Service uses Redis for caching.

The cache key format is:

product:<id>

Example:

product:1

Cache TTL:

10 minutes

PostgreSQL remains the source of truth.

If Redis is unavailable, Product Service can continue using PostgreSQL.

19. Database Development

Each service owns a separate PostgreSQL database.

ecommerce_users
ecommerce_products
ecommerce_orders

Services must not directly access another service's database.

For example:

Order Service
     |
     | REST
     v
Product Service
     |
     v
ecommerce_products
20. Docker Networking

Inside Docker Compose, services communicate using Docker service names.

Examples:

postgres:5432
redis:6379
user-service:8081
product-service:8082
order-service:8083
api-gateway:8080

The API Gateway is the only service exposed to the host in the Docker deployment.

Host access:

http://localhost:8080
21. Development Workflow

Recommended development workflow:

1. Create or update code
        |
        v
2. gofmt
        |
        v
3. go test ./...
        |
        v
4. go vet ./...
        |
        v
5. Docker build
        |
        v
6. Docker Compose E2E test
        |
        v
7. Review git diff
        |
        v
8. Commit
        |
        v
9. Push

Format Go code using:

gofmt -w .
22. Useful Git Commands

Check modified files:

git status

Review changes:

git diff

Review staged changes:

git diff --cached

Stage all changes:

git add .

Commit:

git commit -m "complete phase 1 microservices setup"

Push:

git push origin main
23. Troubleshooting
Check Docker containers
docker compose -f deployments/docker-compose.yml ps
Check service logs
docker compose -f deployments/docker-compose.yml logs -f <service>
Check PostgreSQL
docker compose -f deployments/docker-compose.yml exec postgres \
  psql -U postgres -l
Check Redis
docker compose -f deployments/docker-compose.yml exec redis \
  redis-cli ping

Expected response:

PONG
Check Gateway
curl http://localhost:8080/health
24. Phase 1 Development Scope

Phase 1 provides:

REST-based microservices
API Gateway
User registration and login
JWT authentication
Product CRUD
Redis product caching
Order creation
Order retrieval
Order cancellation
Order ownership authorization
PostgreSQL database-per-service
Docker Compose deployment
Unit tests
Static analysis

Phase 2 will introduce gRPC for selected internal service-to-service
communication and additional distributed-system capabilities.
