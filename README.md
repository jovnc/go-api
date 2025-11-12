# Go API

A simple REST API built with Go, featuring user authentication and management.

## Features

- User registration and authentication
- JWT-based authentication
- PostgreSQL database integration
- Redis integration for caching and session management
- Health check endpoint
- Graceful server shutdown
- Hot reloading in development mode (using Air)

## Prerequisites

- Go 1.25 or higher
- PostgreSQL database
- Redis server
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd go_api
```

2. Install dependencies:

```bash
make deps
# or
go mod download
go mod tidy
```

3. Set up environment variables:
   Create a `.env` file in the root directory with the following variables:

```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
JWT_SECRET_KEY=your-secret-key-here
SERVER_PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Running the Application

### Development Mode (with hot reload)

```bash
make dev
# or
go tool air
```

### Production Mode

```bash
make start
# or
make build
./bin/gobin
```

### Database Migration

```bash
make migrate
```

## Available Make Commands

- `make dev` - Run the application in development mode with hot reload
- `make build` - Build the application binary
- `make start` - Build and start the application
- `make migrate` - Run database migrations
- `make deps` - Download and tidy dependencies
- `make fmt` - Format the code
- `make clean` - Clean build artifacts
- `make stop` - Stop running application

## API Endpoints

### Health Check

- `GET /health` - Check API health status

### User Management

- `POST /users/register` - Register a new user
- `POST /users/login` - Login and get JWT token
- `GET /users/profile` - Get user profile (requires authentication)

## Project Structure

```
go_api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── auth/                # Authentication middleware
│   ├── config/              # Configuration management
│   ├── database/            # Database connection, migrations, and Redis
│   ├── dto/                 # Data transfer objects
│   ├── handlers/            # HTTP handlers
│   ├── models/              # Database models
│   ├── routes/              # Route definitions
│   └── utils/               # Utility functions (JWT, password hashing, etc.)
├── Makefile                 # Build and run commands
└── go.mod                   # Go module dependencies
```

## Technologies Used

- [GORM](https://gorm.io/) - ORM for database operations
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Redis](https://redis.io/) - In-memory data store for caching and session management
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web Tokens for authentication
- [Air](https://github.com/air-verse/air) - Live reload for Go applications
- [godotenv](https://github.com/joho/godotenv) - Environment variable management
- [validator](https://github.com/go-playground/validator) - Input validation
