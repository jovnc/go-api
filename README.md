# Go API

A simple production-ready REST API built with Go, featuring user authentication, blog management, and comprehensive middleware support.

## Features

- User registration and authentication
- JWT-based authentication with token blacklisting
- Blog management (create, read, list, delete)
- Request logging middleware
- Panic recovery middleware
- Rate limiting middleware
- CORS middleware
- Health check endpoint
- Graceful server shutdown
- Hot reloading in development mode (using Air)
- Docker and Docker Compose support
- Github Actions CI/CD pipeline

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

3. Set up environment variables: rename the `.env.example` file to `.env` and fill in the values.

## Running the Application

### Development Mode (with hot reload)

```bash
make dev
```

### Production Mode

```bash
make start
```

### Using Docker Compose

Ensure you have Docker and Docker Compose installed, then run:

```bash
docker-compose up -d
```

## Available Make Commands

- `make dev` - Run the application in development mode with hot reload
- `make build` - Build the application binary
- `make start` - Build and start the application
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
- Add HTTPS/TLS support
- Add request timeout configuration
- Add API documentation (Swagger/OpenAPI)
- Add pagination for list endpoint
- Add request/response structured logging
- Deployment to cloud provider
