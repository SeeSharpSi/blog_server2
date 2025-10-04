# Agent Instructions for Go HTMX Web App

## Build/Test Commands
- **Build**: `go build` or `make build`
- **Run**: `go run .` or `make run`
- **Dev with hot reload**: `make dev` (requires air)
- **Test all**: `go test ./...` or `make test`
- **Test with coverage**: `make test-cover`
- **Lint**: `go vet ./...` or `make lint`
- **Format**: `make fmt`
- **Install tools**: `make install-tools`

## Development Tools
- **Hot Reload**: Air configuration (`.air.toml`) for automatic rebuilding
- **Makefile**: Common development commands in `Makefile`
- **Docker**: Multi-stage Dockerfile and docker-compose.yml
- **Health Check**: `/health` endpoint for monitoring

## Configuration Management
- **Environment Variables**: All configuration is managed via environment variables
- **Environment Support**: Supports development, staging, and production environments via `ENV` variable
- **Configuration File**: Copy `.env.example` to `.env` for local development (automatically loaded)
- **Validation**: Configuration is validated on startup with detailed error messages

## Code Organization
- **Dependency Injection**: Services are injected into handlers for testability
- **Business Logic Separation**: Business logic lives in services, handlers handle HTTP
- **Context Propagation**: Request contexts are passed through the call chain
- **Input Validation**: Dedicated validation package with sanitization utilities
- **Structured Data**: Services return structured data for handlers to use

## Code Style Guidelines
- **Formatting**: Use `gofmt` for Go files, `templ fmt` for .templ files
- **Imports**: Group standard library, third-party, then local imports
- **Naming**: PascalCase for exported identifiers, camelCase for unexported
- **Error Handling**: Return errors, use `errors.Is()` for checking
- **Structs**: Align fields for readability, use meaningful names
- **HTTP Handlers**: Follow standard Go patterns with ResponseWriter and Request
- **Templating**: Use request context for rendering, keep templates clean
- **Sessions**: Use secure random IDs, HTTP-only cookies, 24h expiration
- **Logging**: Use structured logging with slog, include request IDs for tracing
- **Validation**: Use dedicated validation package for input sanitization