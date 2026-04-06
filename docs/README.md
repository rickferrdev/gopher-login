# 🐹 Gopher Login API

![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber_v3-00B7FF?style=for-the-badge&logo=gofiber&logoColor=white)
![Logbull](https://img.shields.io/badge/Logbull-Logging-FF4500?style=for-the-badge&logo=logstash&logoColor=white)
![Uber Fx](https://img.shields.io/badge/Uber_Fx-DI-276DC3?style=for-the-badge&logo=uber&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Clean Architecture](https://img.shields.io/badge/Clean-Architecture-brightgreen?style=for-the-badge)

**Gopher Login** is a high-performance API for authentication and identity management, built with a focus on **Clean Architecture** and **Dependency Injection**. This project was developed as a robust backend-first foundation to support future frontend interfaces, prioritizing decoupling and testability.

> **Disclaimer:** This project is for strictly educational purposes. Although it implements industry standards, a full security audit is recommended before any use in production environments.

---

## 🚀 Core Features

- **Stateless Authentication:** JWT (JSON Web Tokens) implementation for secure and scalable sessions.
- **Security First:** Password hashing using `bcrypt` (via `go-hasher`) and route protection through Guard Middleware.
- **Consumer Management:** Complete registration, login, and profile retrieval (`/me`) flows.
- **Resilience & Performance:** Native rate limiting and configurable context-level timeouts.
- **Observability:** Integration with `slog` and `Logbull` for structured event and error tracking.

---

## 🛠️ Tech Stack

| Component                | Technology                                   |
| :----------------------- | :------------------------------------------- |
| **Runtime**              | Go 1.25.5                                    |
| **Web Framework**        | [Fiber v3](https://docs.gofiber.io/)         |
| **Dependency Injection** | [Uber Fx](https://uber-go.github.io/fx/)     |
| **Persistence (ORM)**    | [Bun](https://bun.uptrace.dev/) (PostgreSQL) |
| **Validation**           | Go-Playground Validator v10                  |
| **Logging**              | Slog + Logbull Adapter                       |

---

## 🏗️ Project Structure

The organization follows **Clean Architecture** principles:

- `cmd/`: Application entry point.
- `internal/api/core/domain/`: Enterprise business entities.
- `internal/api/core/service/`: Business logic and use cases.
- `internal/api/in/rest/`: Input adapters (Fiber Handlers and Middlewares).
- `internal/api/out/database/`: Output adapters (Persistence with Bun).
- `internal/api/platform/`: Cross-cutting tools (Token generation, Validation).

---

## 🚦 Quick Start Guide

### 1. Environment Configuration

The project uses `rickferrdev/dotenv`. Create a `.env` file in the root directory based on `.env.example`:

```env
GOPHER_SERVER_PORT=8080
GOPHER_SERVER_JWT_SECRET=your_super_protected_secret

GOPHER_POSTGRES_URL=postgres://user:pass@localhost:5437/dbname?sslmode=disable
GOPHER_POSTGRES_USER=user
GOPHER_POSTGRES_PASSWORD=pass
GOPHER_POSTGRES_PORT=5437
GOPHER_POSTGRES_DB=dbname

GOPHER_LOGBULL_PROJECT_ID=your_project_id
GOPHER_LOGBULL_HOST=http://localhost:4005
```

### 2\. Infrastructure (Docker)

The repository includes a Compose file to spin up the database and the logging ecosystem. Run the following command:

```bash
docker compose -f docker/compose.yml up -d
```

### 3\. Running the API

```bash
go mod tidy
go run cmd/main.go
```

---

## 🐳 Development with Dev Containers

This repository is ready to use with **VS Code Dev Containers**. When you open the project, VS Code will suggest reopening it in a container, which already includes the configured Go 1.25 environment and necessary Docker extensions.

---

## 📖 API Documentation

### Authentication Endpoints

| Method | Route                   | Description                | Access |
| :----- | :---------------------- | :------------------------- | :----- |
| `POST` | `/api/v1/auth/register` | New user registration      | Public |
| `POST` | `/api/v1/auth/login`    | Login and Token generation | Public |

### User Endpoints

| Method | Route                         | Description         | Access  |
| :----- | :---------------------------- | :------------------ | :------ |
| `GET`  | `/api/v1/consumers/me`        | Logged-in user data | Private |
| `GET`  | `/api/v1/consumers/:username` | Search by username  | Private |

---

## 🔮 Todo

- [ ] Expand Unit and Integration Test coverage (Initial tests set up using GoMock & Testify).
- [ ] Support for Refresh Tokens.
- [ ] Redis integration for Session Caching.
- [ ] **Frontend:** Development of a SPA (React/Next.js) to consume this API.

---

**Developed by [Rickferrdev](https://github.com/rickferrdev)**

---

## 🤝 Contributing

Contributions are what make the open-source community an amazing place to learn, inspire, and create. Any contribution you make is **greatly appreciated**.

## 📄 License

Distributed under the **MIT License**. See the [LICENSE](https://www.google.com/search?q=./LICENSE) file for more details.
