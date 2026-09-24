# E-Commerce Microservices — Database Documentation

## 1. Database Architecture

The application follows a database-per-service approach.

Each microservice owns its own PostgreSQL database.

```text
User Service
     |
     v
ecommerce_users

Product Service
     |
     +------> ecommerce_products
     |
     +------> Redis

Order Service
     |
     v
ecommerce_orders
2. PostgreSQL

PostgreSQL is used as the primary persistent data store.

The Docker deployment uses:

PostgreSQL 16

The following databases are created:

ecommerce_users
ecommerce_products
ecommerce_orders
3. User Database

Database:

ecommerce_users

Owned by:

User Service

The User Service stores user information.

User table

Primary fields:

id
name
email
password_hash
created_at
updated_at

Passwords are never stored as plain text.

The User Service stores a bcrypt password hash.

The API response does not expose password_hash.

4. Product Database

Database:

ecommerce_products

Owned by:

Product Service
Product table

Primary fields:

id
name
description
price
stock
created_at
updated_at

PostgreSQL is the source of truth for product data.

5. Order Database

Database:

ecommerce_orders

Owned by:

Order Service
Order table

Primary fields:

id
user_id
product_id
quantity
total_price
status
created_at
updated_at
Order statuses
PENDING
CONFIRMED
CANCELLED

When an order is created, Order Service retrieves product information from Product Service.

The order total is calculated using:

product price × quantity

Phase 1 does not deduct product stock when an order is created.

This keeps the initial implementation focused on microservice communication without introducing stock-concurrency complexity.

6. Database Ownership
Service	Database	Responsibility
User Service	ecommerce_users	Users and authentication
Product Service	ecommerce_products	Product catalog
Order Service	ecommerce_orders	Orders

Each service is responsible for its own database.

7. Service-to-Service Data Access

Services do not directly query another service's database.

For example, Order Service does not directly access:

ecommerce_products

Instead:

Order Service
     |
     | REST API
     v
Product Service
     |
     v
ecommerce_products

This keeps database ownership within the responsible service.

8. Redis Cache

Redis is used by Product Service as a cache.

Redis is not the primary source of product data.

Cache key

Product cache keys use:

product:<id>

Example:

product:1
Cache flow
Get Product
     |
     v
Check Redis
     |
     +---- Cache Hit ----> Return product
     |
     +---- Cache Miss
              |
              v
         PostgreSQL
              |
              v
        Store in Redis
              |
              v
        Return product

The Product Service cache uses a 10-minute TTL.

9. Redis Failure Handling

Redis is treated as a cache rather than a required persistent dependency.

If Redis is unavailable, Product Service can continue using PostgreSQL.

PostgreSQL remains the source of truth.

10. Docker Database Configuration

Inside Docker Compose, services communicate with PostgreSQL using:

postgres:5432

Redis is accessed using:

redis:6379

The PostgreSQL databases are initialized through:

deployments/postgres/init/

PostgreSQL data is persisted using:

postgres_data

Redis data is persisted using:

redis_data
11. Database Isolation

The architecture avoids sharing one database across all services.

User Service
     |
     +--> ecommerce_users

Product Service
     |
     +--> ecommerce_products

Order Service
     |
     +--> ecommerce_orders

This provides clear ownership boundaries.

12. Phase 1 Database Scope

Phase 1 focuses on:

PostgreSQL persistence
Database-per-service
GORM-based data access
Redis product caching
Service-level database ownership
REST-based service communication

Future phases can introduce:

Distributed consistency
Event-driven updates
Stock reservation
Idempotency
Database scaling
Read replicas
Advanced caching strategies
