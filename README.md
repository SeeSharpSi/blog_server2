# HTMX Quickstart - Go Web Application

A modern, production-ready HTMX quickstart built with Go, featuring clean architecture, comprehensive tooling, and best practices for web development.

## 🚀 Features

### Core Functionality
- **HTMX Integration**: Dynamic web interactions without JavaScript complexity
- **Session Management**: Secure cookie-based sessions with configurable options
- **Health Checks**: Built-in `/health` endpoint for monitoring
- **Error Handling**: Custom 404/500 error pages with structured logging

### Developer Experience
- **Hot Reload**: Automatic rebuilding with Air during development
- **Docker Support**: Multi-stage Docker builds with health checks
- **Makefile**: Convenient commands for common development tasks
- **Structured Logging**: JSON logging with request tracing and context

### Configuration & Deployment
- **Environment Variables**: Flexible configuration via `.env` files
- **Multi-Environment**: Support for development, staging, and production
- **Graceful Shutdown**: Proper signal handling and cleanup
- **Security**: Configurable session security and CORS-ready

### Architecture
- **Clean Architecture**: Dependency injection with service layer separation
- **Input Validation**: Comprehensive validation and sanitization utilities
- **Context Propagation**: Proper request context handling throughout the stack
- **Type Safety**: Strongly typed data structures and interfaces

## 📋 Prerequisites

- Go 1.23 or later
- Docker (optional, for containerized development)
- Make (optional, for using Makefile commands)

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-htmx-web-app
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Install development tools** (optional)
   ```bash
   make install-tools
   ```

## 🚀 Usage

### Development

**Start with hot reload:**
```bash
make dev
```

**Or run directly:**
```bash
make run
```

**Run tests:**
```bash
make test
```

**Format code:**
```bash
make fmt
```

### Production

**Build the application:**
```bash
make build
```

**Run the built binary:**
```bash
./bin/htmx_quickstart
```

### Docker

**Build and run with Docker:**
```bash
make docker-build
make docker-run
```

**Or use docker-compose:**
```bash
docker-compose up
```

## 📁 Project Structure

```
├── config/           # Configuration management
├── handlers/         # HTTP handlers (thin layer)
├── logger/           # Structured logging utilities
├── services/         # Business logic layer
├── session/          # Session management
├── static/           # Static assets (CSS, JS, images)
├── templ/            # HTML templates
├── validation/       # Input validation utilities
├── .air.toml         # Hot reload configuration
├── docker-compose.yml # Docker orchestration
├── Dockerfile        # Container build configuration
├── Makefile          # Development commands
├── .env.example      # Environment configuration template
└── AGENTS.md         # Agent instructions for AI coding assistants
```

## ⚙️ Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and modify as needed:

```bash
# Environment
ENV=development

# Server
SERVER_HOST=localhost
SERVER_PORT=9779
SERVER_ADDRESS=http://localhost

# Session
SESSION_COOKIE_NAME=session_id
SESSION_MAX_AGE=24h
SESSION_SECURE=false
SESSION_HTTP_ONLY=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Database (for future use)
DB_DRIVER=sqlite3
DB_HOST=localhost
DB_NAME=app.db
```

### Environment-Specific Settings

- **Development**: Debug logging, insecure cookies, local database
- **Staging**: Info logging, secure cookies, staging database
- **Production**: Info logging, secure cookies, production database

## 🏗️ Architecture

### Clean Architecture Pattern

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Handlers      │    │   Services      │    │   Repositories  │
│                 │    │                 │    │                 │
│ • HTTP routes   │◄──►│ • Business      │◄──►│ • Data access   │
│ • Request/      │    │   logic         │    │ • Database      │
│   Response      │    │ • Validation    │    │ • External APIs │
│ • Middleware    │    │ • Orchestration │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Key Principles

1. **Dependency Injection**: Services are injected into handlers
2. **Single Responsibility**: Each layer has a specific purpose
3. **Context Propagation**: Request context flows through all layers
4. **Error Handling**: Errors are handled at appropriate layers
5. **Testability**: Clean interfaces enable easy mocking

## 🔍 API Endpoints

- `GET /` - Main index page
- `GET /test` - Test page
- `GET /health` - Health check (JSON response)
- `GET /static/*` - Static file serving

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run tests verbosely
make test-verbose
```

## 📦 Deployment

### Docker Deployment

```bash
# Build production image
docker build -t my-app .

# Run container
docker run -p 9779:8080 --env-file .env my-app
```

### Systemd Service (Linux)

Create `/etc/systemd/system/htmx-quickstart.service`:

```ini
[Unit]
Description=Go HTMX Web Application
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/htmx-quickstart
ExecStart=/opt/htmx-quickstart/bin/htmx_quickstart
Restart=always
Environment=ENV=production

[Install]
WantedBy=multi-user.target
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes with tests
4. Run the test suite: `make test`
5. Format code: `make fmt`
6. Commit your changes: `git commit -am 'Add my feature'`
7. Push to the branch: `git push origin feature/my-feature`
8. Submit a pull request

### Development Guidelines

- Follow Go naming conventions
- Add tests for new functionality
- Update documentation for API changes
- Use structured logging for debugging
- Keep handlers thin, put logic in services

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [HTMX](https://htmx.org/) for the excellent frontend framework
- [templ](https://github.com/a-h/templ) for type-safe HTML templates
- [Air](https://github.com/cosmtrek/air) for hot reloading
- [Go](https://golang.org/) for the amazing programming language

## 📞 Support

If you have questions or need help:

1. Check the [Issues](https://github.com/your-repo/issues) page
2. Review the [AGENTS.md](AGENTS.md) for development guidelines
3. Create a new issue with detailed information

---

**Happy coding!** 🎉