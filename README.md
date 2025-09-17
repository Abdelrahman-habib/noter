# Snippetbox

A secure web application for sharing code snippets, built with Go and MySQL. Features user authentication, session management, and HTTPS support.

## Features

- 📝 Create and share code snippets
- 🔐 User registration and authentication
- 🍪 Secure session management with MySQL store
- 🔒 CSRF protection and bcrypt password hashing
- 🌐 HTTPS/TLS support
- 📱 Responsive web interface
- 🐳 Docker support for easy deployment

## Technology Stack

- **Language**: Go 1.25.1
- **Database**: MySQL 8.0
- **Web Framework**: Standard library `net/http` with custom routing
- **Template Engine**: Go's `html/template`
- **Session Management**: SCS with MySQL store
- **Database Migrations**: Goose
- **Security**: CSRF protection, bcrypt password hashing

## Quick Start with Docker

### Prerequisites

- Docker and Docker Compose installed
- Git

### Development Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd snippetbox
   ```

2. **Start development environment**

   ```bash
   ./docker-run.sh up-dev
   ```

3. **Access the application**
   - Open https://localhost:4444 in your browser
   - Accept the self-signed certificate warning for development

### Production Deployment

```bash
./docker-run.sh up
```

## Manual Setup (Without Docker)

### Prerequisites

- Go 1.25.1 or later
- MySQL 8.0 or later
- Make (optional, for using Makefile commands)

### Installation

1. **Clone and setup**

   ```bash
   git clone <repository-url>
   cd snippetbox
   go mod download
   ```

2. **Database setup**

   ```bash
   # Create database
   mysql -u root -p -e "CREATE DATABASE snippetbox;"

   # Copy environment file
   cp .env.example .env.development
   # Edit .env.development with your database credentials

   # Run migrations
   make migrate-up ENV=development

   # Seed development data (optional)
   make seed-up ENV=development
   ```

3. **Generate TLS certificates** (for HTTPS)

   ```bash
   mkdir tls
   # Generate self-signed certificates for development
   openssl req -x509 -newkey rsa:4096 -keyout tls/key.pem -out tls/cert.pem -days 365 -nodes
   ```

4. **Run the application**
   ```bash
   make run-dev
   ```

## Development

### Available Make Commands

```bash
# Development
make run-dev          # Run in development mode
make build-dev        # Build development binary

# Production
make run-prod         # Run in production mode
make build-prod       # Build production binary

# Testing
make test             # Run all tests
make test-cover       # Run tests with coverage report
make init-test-db     # Initialize test database
make teardown-test-db # Clean up test database

# Database Management
make migrate-up       # Run migrations
make migrate-down     # Rollback migrations
make migrate-reset    # Reset all migrations
make seed-up          # Run seed data (dev/test only)
make seed-down        # Remove seed data

# Code Quality
make lint             # Run go vet
make audit            # Run vet + staticcheck + govulncheck
```

### Docker Development Commands

```bash
# Start services
./docker-run.sh up-dev      # Development mode
./docker-run.sh up          # Production mode

# Development tools
./docker-run.sh test        # Run tests
./docker-run.sh lint        # Run linting
./docker-run.sh audit       # Run security audit
./docker-run.sh tools       # Start tools container

# Logs and debugging
./docker-run.sh logs-web    # View application logs
./docker-run.sh logs-db     # View database logs
./docker-run.sh shell       # Access app container
./docker-run.sh db-shell    # Access MySQL shell

# Cleanup
./docker-run.sh down        # Stop services
./docker-run.sh clean       # Remove all containers and volumes
```

## Project Structure

```
├── cmd/web/              # Application entry point and web server
│   ├── main.go          # Main application bootstrap
│   ├── app.go           # Application struct and server setup
│   ├── config.go        # Configuration parsing
│   ├── handlers.go      # HTTP request handlers
│   ├── middleware.go    # HTTP middleware
│   ├── routes.go        # Route definitions
│   ├── templates.go     # Template rendering logic
│   └── helpers.go       # Helper functions
├── internal/            # Private application packages
│   ├── models/          # Data models and database logic
│   ├── validator/       # Input validation
│   ├── logger/          # Logging utilities
│   └── assert/          # Test assertions
├── db/schema/           # Database schema management
│   ├── migrations/      # Goose migration files
│   └── seed/           # Database seed data
├── ui/                  # User interface assets
│   ├── html/           # HTML templates
│   ├── static/         # CSS, JS, images
│   └── efs.go          # Embedded file system
├── tls/                # TLS certificates
└── bin/                # Built binaries (generated)
```

## Configuration

The application uses environment-specific configuration files:

- `.env.development` - Development settings
- `.env.production` - Production settings
- `.env.test` - Test settings
- `.env` - Fallback configuration

### Key Environment Variables

```bash
# Server Configuration
HOST=localhost
PORT=4444
TLS_CERT=./tls/cert.pem
TLS_KEY=./tls/key.pem

# Database
DB_DSN=user:password@tcp(localhost:3306)/snippetbox

# Application Settings
ENVIROMENT=development
DEBUG=true
```

## Security Features

- **CSRF Protection**: All forms protected against cross-site request forgery
- **Secure Sessions**: HTTP-only, secure cookies with MySQL storage
- **Password Hashing**: bcrypt with appropriate cost factor
- **HTTPS Only**: TLS encryption for all communications
- **Input Validation**: Server-side validation for all user inputs
- **SQL Injection Prevention**: Prepared statements for all database queries

## Testing

### Running Tests

```bash
# Local testing
make test
make test-cover

# Docker testing
./docker-run.sh test
```

### Test Database Setup

```bash
# Initialize test database
make init-test-db

# Clean up after testing
make teardown-test-db
```

## Deployment

### Docker Production Deployment

1. **Prepare environment**

   ```bash
   cp .env.example .env.production
   # Edit .env.production with production values
   ```

2. **Deploy**
   ```bash
   ./docker-run.sh up
   ```

### Manual Production Deployment

1. **Build the application**

   ```bash
   make build-prod
   ```

2. **Setup database**

   ```bash
   make migrate-up ENV=production
   ```

3. **Configure TLS certificates**

   - Replace self-signed certificates in `tls/` with proper certificates
   - Or configure reverse proxy (nginx/Apache) for TLS termination

4. **Run the application**
   ```bash
   make run-prod
   ```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test` or `./docker-run.sh test`)
5. Run linting (`make audit` or `./docker-run.sh audit`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:

- Create an issue in the GitHub repository
- Check the documentation in the `docs/` directory
- Review the code comments for implementation details
