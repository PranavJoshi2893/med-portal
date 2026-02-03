# Med Portal

A REST API for user management with JWT authentication, built in Go.

## Tech Stack

- **Go** – Backend
- **Chi** – HTTP router
- **PostgreSQL** – Database
- **JWT** – Access & refresh tokens

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL (via Docker, or local on port 5432)

## Setup

### 1. Environment

Create a `.env` file from `example.env`:

```bash
cp example.env .env
```

Configure at least:

- `PORT` – Server port (e.g. `:3000`)
- `PEPPER` – Password hashing pepper
- `ACCESS_TOKEN_KEY` – JWT access token secret
- `REFRESH_TOKEN_KEY` – JWT refresh token secret
- `DATABASE_URL` – PostgreSQL connection string

### 2. Database

Start PostgreSQL:

```bash
docker-compose up -d
```

Run migrations:

```bash
make migrate-up-all
```

### 3. Build & Run

```bash
make build
make run
```

The API runs at `http://localhost:3000` (or your configured port).

## API Reference

Base URL: `/api/v1`

### Auth (public)

| Method | Endpoint        | Description                    |
|--------|-----------------|--------------------------------|
| POST   | `/auth/register`| Register a new user            |
| POST   | `/auth/login`   | Login, returns access token    |

### Auth (refresh token required)

| Method | Endpoint       | Description                    |
|--------|----------------|--------------------------------|
| POST   | `/auth/refresh`| Issue new access token        |
| POST   | `/auth/logout` | Revoke refresh token           |

### Users (access token required)

| Method | Endpoint    | Description                    |
|--------|-------------|--------------------------------|
| GET    | `/users/`   | List users (paginated)        |
| GET    | `/users/{id}`| Get user by ID                |
| PATCH  | `/users/{id}`| Update user                   |
| DELETE | `/users/{id}`| Soft delete user              |

### Pagination

`GET /users/` supports:

- `page` – Page number (default: 1)
- `limit` – Page size: **5**, **10**, or **100** (default: 10)

Example: `GET /users/?page=1&limit=10`

Response:

```json
{
  "message": "success",
  "data": {
    "items": [...],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    }
  }
}
```

### Authentication

- **Access token**: Send in `Authorization: Bearer <token>` for `/users/*`
- **Refresh token**: Stored in HTTP-only cookie; used for `/auth/refresh` and `/auth/logout`

## Makefile

| Command           | Description              |
|-------------------|--------------------------|
| `make build`      | Build the binary         |
| `make run`        | Build and run            |
| `make clean`      | Remove build artifacts   |
| `make migrate-up` | Run one migration up     |
| `make migrate-up-all` | Run all migrations up |
| `make migrate-down`   | Run one migration down   |

## Scripts

`scripts/requests.sh` – Example flow: register, login, list users, update, refresh, logout.

```bash
./scripts/requests.sh
```

Requires `jq` and a running server.
