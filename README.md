# Health Check API

A simple, production-ready health check API server built with Go and Gin framework, following clean architecture principles.

## Features

- **Clean Architecture**: Separation of concerns with distinct layers (domain, appl[engine.go](../bookmark-management/internal/api/engine.go)ication, infrastructure)
- **Environment Configuration**: Configurable via environment variables
- **Auto-generated UUID**: Automatically generates instance ID if not provided
- **Comprehensive Testing**: Includes unit tests and integration tests
- **API Documentation**: Swagger/OpenAPI documentation included
- **Makefile**: Convenient commands for building, running, and testing

## Architecture

This project follows clean architecture principles with the following structure:

```
health-check/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # Infrastructure layer (HTTP server, routing)
│   │   ├── api.go
│   │   └── config.go
│   ├── handler/                 # HTTP handlers (presentation layer)
│   │   ├── healthcheck.go
│   │   └── healthcheck_test.go
│   ├── service/                 # Business logic layer
│   │   ├── healthcheck.go
│   │   └── healthcheck_test.go
│   ├── model/                   # Domain entities
│   │   └── healthcheck.go
│   └── intergration_test/       # Integration tests
│       └── healthcheck_test.go
├── docs/                        # Swagger documentation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Layers

- **Domain Layer (`model/`)**: Contains business entities and core business rules
- **Application Layer (`service/`)**: Contains use cases and business logic
- **Infrastructure Layer (`api/`, `handler/`)**: Handles HTTP, routing, and external concerns

## Installation

### Prerequisites

- Go 1.25 or higher
- Make (optional, for using Makefile commands)

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd health-check
```

2. Install dependencies:
```bash
go mod download
```

3. Install swag for Swagger documentation (optional):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Configuration

The application is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVICE_NAME` | Name of the service | `health-check-service` |
| `INSTANCE_ID` | Unique instance identifier | Auto-generated UUID if empty |
| `PORT` | Port to run the server on | `8080` |

### Example Configuration

```bash
# Using environment variables
export SERVICE_NAME=my-service
export INSTANCE_ID=123e4567-e89b-12d3-a456-426614174000
export PORT=8080

# Or run with inline variables
SERVICE_NAME=my-service INSTANCE_ID=test-123 PORT=8080 go run cmd/api/main.go
```

## Running the Application

### Using Makefile

```bash
# Run the application
make run

# Build the application
make build

# Run the built binary
./build/health-check
```

### Using Go directly

```bash
go run cmd/api/main.go
```

The server will start on the configured port (default: 8080).

## API Documentation

### Swagger UI

Once the server is running, access the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

### Health Check Endpoint

**GET /health-check**

Returns the health status of the service including Redis connection check.

**Response (Success):**
```json
{
  "message": "OK",
  "service_name": "health-check-service",
  "instance_id": "18cf5eaf-f401-45bf-a827-265b9c6c9333"
}
```

**Response (Redis Error):**
```json
{
  "error": "Redis connection failed"
}
```

**Status Code:** 200 (success), 500 (Redis connection failed)

### URL Shortening Endpoint

**POST /v1/links/shorten**

Creates a short code for a given URL.

**Request:**
```json
{
  "url": "https://example.com",
  "exp": 604800
}
```

**Response:**
```json
{
  "code": "abc1234",
  "message": "Shorten URL generated successfully!"
}
```

**Status Code:** 200 (success), 400 (invalid request), 500 (server error)

**Parameters:**
- `url` (string, required): The URL to shorten
- `exp` (int, required): Expiration time in seconds

## Testing

### Run All Tests

```bash
# Using Makefile
make test

# Using Go
go test -v ./...
```

### Run Unit Tests Only

```bash
# Using Makefile
make test-unit

# Using Go
go test -v ./internal/... -run "^((?!integration).)*$$"
```

### Run Integration Tests Only

```bash
# Using Makefile
make test-integration

# Using Go
go test -v ./internal/intergration_test/...
```

### Test Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make all` | Build the application |
| `make build` | Build the application |
| `make run` | Run the application |
| `make test` | Run all tests |
| `make test-integration` | Run integration tests only |
| `make swagger` | Generate Swagger documentation |
## Deployment

This application supports containerized deployment using Docker Compose. For more detailed instructions on building images, routing, and scaling services, please see the [Deployment & Infrastructure Guide](file:///d:/CODE/Golang/golang-dev/deployment/README.md).

### VM Deployment

The application is deployed on a production VM with the IP address `103.75.183.118`. You can access and verify the live environment via:

* **Frontend UI**: [http://103.75.183.118/](http://103.75.183.118/)
* **Health Check Endpoint**: [http://103.75.183.118/api/bookmark_service/health-check](http://103.75.183.118/api/bookmark_service/health-check)
* **Swagger API Documentation**: [http://103.75.183.118/api/bookmark_service/swagger/index.html](http://103.75.183.118/api/bookmark_service/swagger/index.html)
* **Swagger Health Check Detail**: [http://103.75.183.118/api/bookmark_service/swagger/index.html#/health-check/get_health_check](http://103.75.183.118/api/bookmark_service/swagger/index.html#/health-check/get_health_check)


## Development

### Adding New Endpoints

1. Define the model in `internal/model/`
2. Create the service logic in `internal/service/`
3. Create the handler in `internal/handler/`
4. Add the route in `internal/api/api.go`
5. Add Swagger annotations to the handler
6. Regenerate Swagger docs with `make swagger`

### Code Style

- Follow Go best practices and idiomatic Go
- Use `make fmt` to format code before committing
- Run `make lint` to check for code issues

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [envconfig](https://github.com/kelseyhightower/envconfig) - Environment variable configuration
- [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger) - Swagger middleware for Gin
- [testify](https://github.com/stretchr/testify) - Testing toolkit

## License

This project is licensed under the Apache 2.0 License.
