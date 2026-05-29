# auth-server

A Go backend service providing authentication (signup, login) and profile management via a REST API, built with [Fiber](https://github.com/gofiber/fiber), [GORM](https://gorm.io/), and PostgreSQL.

## Tech Stack

- **Language:** Go 1.20
- **Web Framework:** Fiber v2
- **ORM:** GORM
- **Database:** PostgreSQL
- **Auth:** JWT (golang-jwt)
- **Password Hashing:** bcrypt

## API Endpoints

| Method | Endpoint              | Auth Required | Description           |
|--------|-----------------------|---------------|-----------------------|
| GET    | `/api/`               | No            | Root (health check)   |
| POST   | `/api/signup`         | No            | Register a new user   |
| POST   | `/api/login`          | No            | Login and get JWT     |
| GET    | `/api/profile`        | Yes           | Get authenticated user profile |

### Request / Response Examples

**Signup**

```json
POST /api/signup
{
  "username": "johndoe",
  "password": "securepass",
  "email": "john@example.com",
  "full_name": "John Doe"
}
```

**Login**

```json
POST /api/login
{
  "username": "johndoe",
  "password": "securepass"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Profile**

```
GET /api/profile
Authorization: Bearer <token>
```

Response:

```json
{
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_admin": false
  }
}
```

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL

### Installation

```bash
# Clone the repo
git clone <repo-url>
cd auth-server

# Install dependencies
go mod download
```

### Environment Setup

Create a `.env` file in the project root:

```env
SERVER_PORT=3000
SECRET=your-super-secret-jwt-key-change-this

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_db

STATIC_ROOT=./static
MEDIA_ROOT=./media
```

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:3000`.

## Project Structure

```
auth-server/
├── config/          # Environment configuration
├── database/        # Database connection and setup
├── handler/         # HTTP request handlers
├── middleware/      # JWT authentication middleware
├── model/          # GORM data models
├── router/         # Route definitions
├── main.go         # Application entry point
├── go.mod
└── .env            # Local environment variables (not committed)
```

## Database Schema

The `users` table is auto-migrated by GORM:

| Column    | Type         | Constraints           |
|-----------|-------------|-----------------------|
| id        | bigint       | primary key, auto     |
| username  | varchar      | primary key, unique   |
| password  | varchar      | not null (bcrypt)     |
| email     | varchar      | unique                |
| full_name | text         |                       |
| is_admin  | boolean      | default: false        |
| created_at| timestamp    | auto                  |
| updated_at| timestamp    | auto                  |
| deleted_at| timestamp    | auto (soft delete)    |

## License

MIT
