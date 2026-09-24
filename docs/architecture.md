# E-Commerce Microservices — Architecture

## 1. Overview

This project is a learning and portfolio-oriented e-commerce backend built using
Golang and a microservices architecture.

The project is being developed in multiple phases.

### Phase 1

Phase 1 uses:

- Golang
- REST APIs
- PostgreSQL
- Redis
- API Gateway
- JWT authentication
- Docker
- Docker Compose

### Phase 2

Phase 2 will introduce:

- gRPC
- Internal service-to-service communication using gRPC
- Additional distributed-system concepts

---

## 2. High-Level Architecture

```text
                         Client
                           |
                           | HTTP
                           v
                  +-------------------+
                  |    API Gateway    |
                  |       :8080       |
                  +---------+---------+
                            |
             +--------------+--------------+
             |              |              |
             v              v              v
      +-----------+  +-----------+  +-----------+
      |   User    |  |  Product  |  |   Order   |
      |  Service  |  |  Service  |  |  Service  |
      |   :8081   |  |   :8082   |  |   :8083   |
      +-----+-----+  +-----+-----+  +-----+-----+
            |              |              |
            v              v              v
      +-----------+  +-----------+  +-----------+
      | PostgreSQL|  | PostgreSQL|  | PostgreSQL|
      |   Users   |  | Products  |  |  Orders   |
      +-----------+  +-----------+  +-----------+

                         Product Service
                              |
                              v
                           Redis