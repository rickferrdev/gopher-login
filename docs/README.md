````md
# Gopher Login

Gopher Login is a simple authentication API built with Go, designed mainly for study purposes and as a reusable foundation for other backend projects.

The project implements a minimal login and registration flow using a clean project structure inspired by Ports and Adapters architecture. It is intentionally small, but organized enough to be used as a learning reference or as a complement to larger applications.

## Purpose

The main goal of this project is to provide a practical example of how to structure a Go backend application with authentication, dependency injection, HTTP routing, MongoDB persistence, password hashing, and JWT token generation.

It can be used as:

- A study project for Go backend development
- A reference for organizing APIs with modular architecture
- A starting point for authentication in other projects
- A small companion service for applications that need basic user login

## Features

- User registration
- User login
- Password hashing with bcrypt
- JWT token generation
- MongoDB user persistence
- Unique email index
- Environment-based configuration
- HTTP API using Fiber
- Dependency injection with Uber Fx
- Centralized error handling
- Docker support
- Basic request collections for manual API testing

## Tech Stack

- Go
- Fiber v3
- MongoDB
- MongoDB Go Driver
- Uber Fx
- JWT
- bcrypt
- Docker
- Bruno / HTTP files for request testing

## Project Structure

```txt
.
├── cmd/
│   └── api/
│       └── main.go
├── docker/
│   ├── docker-compose.yml
│   └── Dockerfile
├── docs/
│   └── README.md
├── internal/
│   ├── config/
│   ├── core/
│   ├── inbound/
│   ├── infra/
│   ├── outbound/
│   └── platform/
├── .env.example
├── Makefile
├── go.mod
└── LICENSE
````

The application is grouped by responsibility. The main modules include configuration, infrastructure, platform services, outbound adapters, application services, and inbound HTTP handlers.

## Architecture Overview

This project follows a simple modular structure inspired by Ports and Adapters.

### Core

The `core` package contains the application domain, service contracts, storage ports, platform ports, and application services.

The authentication service exposes two main operations:

* `Register`
* `Login`

These operations are defined through the `Auth` service interface.

### Inbound

The `inbound` layer contains entry points into the application, such as HTTP routes and handlers.

Currently, the REST API exposes authentication endpoints for login and registration.

### Outbound

The `outbound` layer contains adapters that communicate with external systems.

In this project, MongoDB is used as the persistence layer for users. The user storage adapter implements operations for finding a user by email and inserting a new user.

### Platform

The `platform` layer contains reusable technical services that are not business-specific, such as password hashing and JWT token generation.

Password hashing is implemented with bcrypt.

JWT generation and validation are implemented using `github.com/golang-jwt/jwt/v5`.

## API Endpoints

The API is mounted under:

```txt
/api/v1
```

### Register User

```http
POST /api/v1/auth/register
Content-Type: application/json
```

#### Request Body

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure123"
}
```

#### Response

```json
{
  "id": "user_id",
  "created_at": "2026-01-01T00:00:00Z"
}
```

The registration endpoint creates a new user, hashes the password, stores the user in MongoDB, and returns the created user ID.

### Login User

```http
POST /api/v1/auth/login
Content-Type: application/json
```

#### Request Body

```json
{
  "email": "john@example.com",
  "password": "secure123"
}
```

#### Response

```json
{
  "token": "jwt_token",
  "created_at": "2026-01-01T00:00:00Z"
}
```

The login endpoint validates the user credentials and returns a JWT token when authentication succeeds.

## Validation Rules

The request DTOs include validation tags:

```go
Email    string `json:"email" validate:"email"`
Password string `json:"password" validate:"min=8"`
Username string `json:"username" validate:"min=4,max=16"`
```

This means:

* `email` must be a valid email address
* `password` must have at least 8 characters
* `username` must have between 4 and 16 characters



## Environment Variables

Create a `.env` file based on `.env.example`:

```env
SERVER_PORT=8080
DATABASE_URI=mongodb://mongo:27017
JWT_SECRET=change-me-in-production
JWT_EXPIRES_IN=24h
```

The application reads its configuration from environment variables, including the server port, MongoDB URI, JWT secret, and JWT expiration time.

## Running Locally

### Requirements

* Go
* MongoDB
* Docker and Docker Compose, optional but recommended

### Using Go directly

```bash
go run ./cmd/api
```

Or using the Makefile:

```bash
make run
```

### Using Docker Compose

```bash
docker compose -f docker/docker-compose.yml up --build
```

The Docker Compose setup starts both the API and a MongoDB container. The API is exposed on port `8080`.

## Available Make Commands

```bash
make run
make test
make tidy
make fmt
make lint
make build
```

These commands are defined in the project Makefile and provide shortcuts for common development tasks.

## Database

The project uses MongoDB as its database.

The default database name is:

```txt
gopher
```

The user collection name is:

```txt
users
```



A unique index is created for the `email` field to prevent duplicated users.

## Authentication Flow

### Registration

1. The client sends `username`, `email`, and `password`.
2. The service validates required fields.
3. The email is normalized.
4. The password is hashed using bcrypt.
5. A new user ID is generated.
6. The user is stored in MongoDB.
7. The API returns the created user ID.

### Login

1. The client sends `email` and `password`.
2. The service validates required fields.
3. The email is normalized.
4. The user is searched by email.
5. The password is compared with the stored hash.
6. A JWT token is generated.
7. The API returns the token.

## Error Handling

The project uses a centralized error structure with internal error codes, HTTP status codes, and user-facing messages.

Example error format:

```json
{
  "code": "INVALID_CREDENTIALS",
  "status": 401,
  "message": "invalid credentials"
}
```

The custom error system maps application errors to proper HTTP responses.

## Study Goals

This project is useful for studying:

* Go project organization
* Clean separation between layers
* Dependency injection with Uber Fx
* HTTP APIs with Fiber
* MongoDB integration
* Password hashing
* JWT authentication
* Error mapping
* Docker-based development
* Basic authentication service design

## Possible Improvements

This project is intentionally simple, but it can be expanded with:

* Refresh tokens
* Email verification
* Password reset flow
* Role-based authorization
* Protected routes examples
* Unit and integration tests
* OpenAPI documentation
* Observability with logs, metrics, and traces
* CI/CD pipeline
* Rate limit customization
* Docker Compose profiles for development and production

## License

This project is licensed under the MIT License.
