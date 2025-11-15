# Go API

A simple production-ready REST API built with Go, featuring user authentication, blog management, and comprehensive middleware support.

## Features

- User registration and authentication
- JWT-based authentication with token blacklisting
- Blog management (create, read, list, delete)
- PostgreSQL database integration
- Redis integration for caching and session management
- Request logging middleware
- Panic recovery middleware
- Rate limiting middleware
- Health check endpoint
- Graceful server shutdown
- Hot reloading in development mode (using Air)

## Prerequisites

- Go 1.25 or higher
- PostgreSQL database
- Redis server
- Make (for running Makefile commands)

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd go_api
```

2. Install dependencies:

```bash
make deps
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
RATE_LIMIT=100
```

## Running the Application

### Development Mode (with hot reload)

```bash
make dev
```

### Production Mode

```bash
make start
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
- `POST /users/logout` - Logout and blacklist token (requires authentication)
- `GET /users/profile` - Get user profile (requires authentication)
- `GET /users/` - List all users (requires authentication)

### Blog Management

- `POST /blogs/` - Create a new blog post (requires authentication)
- `GET /blogs/{id}` - Get a blog post by ID
- `GET /blogs/` - List all blog posts
- `DELETE /blogs/{id}` - Delete a blog post (requires authentication, owner only)

## Project Structure

```
go_api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/
│   │   ├── dto/             # Data transfer objects
│   │   ├── handler/         # HTTP handlers
│   │   ├── model/           # Database models
│   │   └── route/           # Route definitions
│   ├── config/              # Configuration management
│   ├── middleware/          # Middlewares
│   ├── storage/             # Database and Redis connections
│   └── util/                # Utility functions
├── http/                    # HTTP test files
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
- [httprate](https://github.com/go-chi/httprate) - Rate limiting middleware

## Planned Enhancements

- Add comprehensive test suite (unit tests, integration tests, etc.)
- Add CORS middleware
- Add HTTPS/TLS support
- Add request timeout configuration
- Add healthcheck for database and Redis
- Add API documentation (Swagger/OpenAPI)
- Add pagination for list endpoint
- Add request/response structured logging
- Add Docker support
- Add CI/CD pipeline
- Deployment to cloud provider
