# 🐹 Gopher Login API

![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber_v3-00B7FF?style=for-the-badge&logo=gofiber&logoColor=white)
![Logbull](https://img.shields.io/badge/Logbull-Logging-FF4500?style=for-the-badge&logo=logstash&logoColor=white)
![Uber Fx](https://img.shields.io/badge/Uber_Fx-DI-276DC3?style=for-the-badge&logo=uber&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Clean Architecture](https://img.shields.io/badge/Clean-Architecture-brightgreen?style=for-the-badge)

O **Gopher Login** é uma API de alto desempenho para autenticação e gestão de identidade, construída com foco em **Clean Architecture** e **Injeção de Dependência**. Este projeto foi desenvolvido como uma fundação robusta focada no backend para suportar futuras interfaces frontend, priorizando o desacoplamento e a testabilidade.

> **Aviso:** Este projeto tem fins estritamente educacionais. Embora implemente padrões da indústria, uma auditoria de segurança completa é recomendada antes de qualquer uso em ambientes de produção.

---

## 🚀 Funcionalidades Core

- **Autenticação Stateless:** Implementação de JWT (JSON Web Tokens) para sessões seguras e escaláveis.
- **Segurança em Primeiro Lugar:** Hashing de senhas usando `bcrypt` (via `go-hasher`) e proteção de rotas através do Middleware de Guard.
- **Gestão de Usuários:** Fluxos completos de registro, login e busca de perfil (`/me`).
- **Resiliência e Performance:** Rate limiting nativo e timeouts configuráveis em nível de contexto.
- **Observabilidade:** Integração com `slog` e `Logbull` para rastreamento estruturado de eventos e erros.

---

## 🛠️ Stack Tecnológica

| Componente                 | Tecnologia                                   |
| :------------------------- | :------------------------------------------- |
| **Ambiente de Execução**   | Go 1.25.5                                    |
| **Framework Web**          | [Fiber v3](https://docs.gofiber.io/)         |
| **Injeção de Dependência** | [Uber Fx](https://uber-go.github.io/fx/)     |
| **Persistência (ORM)**     | [Bun](https://bun.uptrace.dev/) (PostgreSQL) |
| **Validação**              | Go-Playground Validator v10                  |
| **Logs**                   | Slog + Adaptador Logbull                     |

---

## 🏗️ Estrutura do Projeto

A organização segue os princípios da **Clean Architecture**:

- `cmd/`: Ponto de entrada da aplicação.
- `internal/api/core/domain/`: Entidades de negócio.
- `internal/api/core/service/`: Lógica de negócio e casos de uso.
- `internal/api/in/rest/`: Adaptadores de entrada (Handlers e Middlewares do Fiber).
- `internal/api/out/database/`: Adaptadores de saída (Persistência com Bun).
- `internal/api/platform/`: Ferramentas transversais (Geração de Token, Validação).

---

## 🚦 Guia de Início Rápido

### 1. Configuração de Ambiente

O projeto usa o pacote `rickferrdev/dotenv`. Crie um arquivo `.env` no diretório raiz baseando-se no `.env.example`:

```env
GOPHER_SERVER_PORT=8080
GOPHER_SERVER_JWT_SECRET=seu_segredo_super_protegido

GOPHER_POSTGRES_URL=postgres://user:pass@localhost:5437/dbname?sslmode=disable
GOPHER_POSTGRES_USER=user
GOPHER_POSTGRES_PASSWORD=pass
GOPHER_POSTGRES_PORT=5437
GOPHER_POSTGRES_DB=dbname

GOPHER_LOGBULL_PROJECT_ID=seu_projeto_id
GOPHER_LOGBULL_HOST=http://localhost:4005
```

### 2\. Infraestrutura (Docker)

O repositório inclui um arquivo Compose para subir o banco de dados e o ecossistema de logs. Execute o seguinte comando:

```bash
docker compose -f docker/compose.yml up -d
```

### 3\. Executando a API

```bash
go mod tidy
go run cmd/main.go
```

---

## 🐳 Desenvolvimento com Dev Containers

Este repositório está pronto para uso com **VS Code Dev Containers**. Ao abrir o projeto, o VS Code sugerirá a reabertura em um container, que já inclui o ambiente Go 1.25 configurado e as extensões necessárias para o Docker.

---

## 📖 Documentação da API

### Endpoints de Autenticação

| Método | Rota                    | Descrição                | Acesso  |
| :----- | :---------------------- | :----------------------- | :------ |
| `POST` | `/api/v1/auth/register` | Registro de novo usuário | Público |
| `POST` | `/api/v1/auth/login`    | Login e geração de Token | Público |

### Endpoints de Usuários

| Método | Rota                          | Descrição                 | Acesso  |
| :----- | :---------------------------- | :------------------------ | :------ |
| `GET`  | `/api/v1/consumers/me`        | Dados do usuário logado   | Privado |
| `GET`  | `/api/v1/consumers/:username` | Busca por nome de usuário | Privado |

---

**Desenvolvido por [Rickferrdev](https://github.com/rickferrdev)**

---

## 🤝 Contribuindo

As contribuições são o que tornam a comunidade open-source um lugar incrível para aprender, inspirar e criar. Qualquer contribuição que você fizer será **muito apreciada**.

## 📄 Licença

Distribuído sob a **Licença MIT**. Veja o arquivo [LICENSE](https://www.google.com/search?q=./LICENSE) para mais detalhes.
