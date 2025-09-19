# Noter

A secure web application for sharing Notes, Qoutes, and any text in general, built with Go and MySQL. Features user authentication, session management, and HTTPS support.

## Features

- 📝 Create and share text notes
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
   cd noter
   ```

2. **Start development environment**

   ```bash
   ./scripts/docker-run.sh dev
   ```

3. **Run database migrations** (first time only)

   ```bash
   ./scripts/docker-run.sh migrate
   ```

4. **Access the application**
   - Open https://localhost:4444 in your browser
   - Accept the self-signed certificate warning for development
   - MySQL is available on localhost:3307

### Production Deployment

```bash
./docker-run.sh prod
```

## Docker Setup Details

### Architecture

The Docker setup uses a multi-container architecture:

- **MySQL Container**: Database server with automatic initialization
- **Web Container**: Go application with live reload for development
- **Migration Container**: Runs Goose migrations (on-demand)

### Services

#### MySQL Service (`noter-mysql`)

- **Image**: `mysql:8.0`
- **Port**: `3307:3306` (mapped to avoid conflicts with local MySQL)
- **Data**: Persistent volume `mysql_data`
- **Initialization**: Automatically creates databases and users via `db/init/`
- **Health Check**: Ensures MySQL is ready before other services start

#### Web Service (`noter-web`)

- **Build**: Custom Go application image
- **Port**: `4444:4444` (HTTPS)
- **Volumes**:
  - Source code mounted for live reload
  - Go module cache for faster builds
- **Dependencies**: Waits for MySQL to be healthy
- **Environment**: Configurable via environment variables

#### Migration Service (`noter-migrate`)

- **Profile**: `migrate` (runs on-demand)
- **Purpose**: Runs Goose migrations and seed data
- **Dependencies**: Waits for MySQL to be healthy
- **Command**: `goose -dir db/schema/migrations up && goose -dir db/schema/seed -no-versioning up`

### Database Initialization

The MySQL container automatically initializes with:

1. **Database Creation**: `noter` and `noter_test` databases
2. **User Creation**: Three users with appropriate permissions:
   - `noter_admin`: Full privileges on `noter` database (for migrations)
   - `noter_web`: Limited privileges on `noter` database (for application)
   - `noter_test_web`: Full privileges on `noter_test` database (for tests)
3. **Migration Execution**: Goose runs migrations to create tables
4. **Seed Data**: Sample data inserted for development

### Environment Configuration

#### Development Environment (`dev.env`)

```env
# Server Configuration
HOST=localhost
PORT=4000
ENVIROMENT=development

# Database Configuration
GOOSE_DRIVER=mysql
GOOSE_DBSTRING=noter_admin:admin@tcp(mysql:3306)/noter
GOOSE_MIGRATION_DIR=./db/schema/migrations
GOOSE_SEED_DIR=./db/schema/seed
GOOSE_TABLE=noter.goose_migrations
DB_DSN=noter_web:pass@tcp(mysql:3306)/noter
TEST_DB_DSN=noter_test_web:test_pass@tcp(mysql:3306)/noter_test

# Application Settings
DEBUG=true
LOG_LEVEL=debug
TLS_CERT=./tls/cert.pem
TLS_KEY=./tls/key.pem

# Docker-specific configuration
MYSQL_ROOT_PASSWORD=dev_root_password_change_me
```

#### Security Considerations

- **MySQL Root Password**: Configurable via `MYSQL_ROOT_PASSWORD` environment variable
- **Default Password**: `dev_root_password_change_me` (change for production)
- **Database Users**: Separate users with minimal required permissions
- **Network Isolation**: All services run in isolated Docker network

### Docker Commands

#### Service Management

```bash
# Start development environment
./docker-run.sh dev

# Start production environment
./docker-run.sh prod

# Stop all services
./docker-run.sh down

# View service status
docker-compose ps
```

#### Database Management

```bash
# Run migrations
./docker-run.sh migrate          # Run migrations (dev)
./docker-run.sh migrate-dev      # Run migrations (dev)
./docker-run.sh migrate-prod     # Run migrations (prod)

# Migration operations
./docker-run.sh migrate-up       # Run migrations up (dev)
./docker-run.sh migrate-down     # Rollback migrations (dev)
./docker-run.sh migrate-reset    # Reset all migrations (dev)
./docker-run.sh migrate-status   # Check migration status (dev)

# Database access
./docker-run.sh db-shell         # Access MySQL shell
./docker-run.sh logs-db          # View database logs
```

#### Development Tools

```bash
# Run tests
./docker-run.sh test

# Run linting
./docker-run.sh lint

# Run security audit
./docker-run.sh audit

# Access tools container
./docker-run.sh tools
```

#### Debugging

```bash
# View application logs
./docker-run.sh logs-web

# Access application container
./docker-run.sh shell

# Clean up everything
./docker-run.sh clean
```

### Troubleshooting

#### Common Issues

1. **Port Conflicts**

   - MySQL uses port 3307 to avoid conflicts with local MySQL
   - Web app uses port 4444 for HTTPS

2. **Database Connection Issues**

   - Ensure MySQL container is healthy before running migrations
   - Check that database users have correct permissions

3. **Migration Failures**

   - Run `./docker-run.sh migrate` to apply migrations
   - Check migration logs for specific errors

4. **Certificate Issues**
   - Self-signed certificates are used for development
   - Accept the security warning in your browser
   - For production, replace with proper certificates

#### Reset Everything

```bash
# Stop and remove all containers, networks, and volumes
./docker-run.sh clean

# Start fresh
./docker-run.sh dev
./docker-run.sh migrate
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
   cd noter
   go mod download
   ```

2. **Database setup**

   ```bash
   # Create database
   mysql -u root -p -e "CREATE DATABASE noter;"

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
./scripts/docker-run.sh dev         # Development mode
./scripts/docker-run.sh prod        # Production mode

# Database management
./scripts/docker-run.sh migrate         # Run database migrations (dev)
./scripts/docker-run.sh migrate-dev     # Run database migrations (dev)
./scripts/docker-run.sh migrate-prod    # Run database migrations (prod)
./scripts/docker-run.sh migrate-up      # Run migrations up (dev)
./scripts/docker-run.sh migrate-down    # Run migrations down (dev)
./scripts/docker-run.sh migrate-reset   # Reset all migrations (dev)
./scripts/docker-run.sh migrate-status  # Check migration status (dev)

# Development tools
./scripts/docker-run.sh test        # Run tests
./scripts/docker-run.sh lint        # Run linting
./scripts/docker-run.sh audit       # Run security audit
./scripts/docker-run.sh tools       # Start tools container

# Logs and debugging
./scripts/docker-run.sh logs-web    # View application logs
./scripts/docker-run.sh logs-mysql  # View database logs
./scripts/docker-run.sh shell       # Access app container
./scripts/docker-run.sh db-shell    # Access MySQL shell

# Cleanup
./scripts/docker-run.sh down        # Stop services
./scripts/docker-run.sh clean       # Remove all containers and volumes
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
├── db/                  # Database configuration
│   ├── schema/          # Database schema management
│   │   ├── migrations/  # Goose migration files
│   │   └── seed/        # Database seed data
│   └── init/            # MySQL initialization scripts
│       └── 01-create-databases-and-users.sql
├── ui/                  # User interface assets
│   ├── html/           # HTML templates
│   ├── static/         # CSS, JS, images
│   └── efs.go          # Embedded file system
├── scripts/            # Helper scripts
│   ├── docker-run.sh  # Docker helper script
│   └── entrypoint.sh  # Container entrypoint script
├── tls/                # TLS certificates
├── bin/                # Built binaries (generated)
├── docker-compose.yml  # Docker Compose configuration
├── Dockerfile          # Go application Docker image
├── dev.env            # Development environment variables
├── prod.env           # Production environment variables
├── env.template       # Environment variables template
├── makefile           # Make commands for development
```

## Configuration

The application uses environment-specific configuration files:

- `dev.env` - Development settings (Docker)
- `env.template` - Template for creating environment files
- `.env.development` - Development settings (manual setup)
- `.env.production` - Production settings
- `.env.test` - Test settings
- `.env` - Fallback configuration

### Docker Environment Variables

The Docker setup uses `dev.env` for development configuration:

```bash
# Server Configuration
HOST=localhost
PORT=4000
ENVIROMENT=development

# Database Configuration
GOOSE_DRIVER=mysql
GOOSE_DBSTRING=noter_admin:admin@tcp(mysql:3306)/noter
GOOSE_MIGRATION_DIR=./db/schema/migrations
GOOSE_SEED_DIR=./db/schema/seed
GOOSE_TABLE=noter.goose_migrations
DB_DSN=noter_web:pass@tcp(mysql:3306)/noter
TEST_DB_DSN=noter_test_web:test_pass@tcp(mysql:3306)/noter_test

# Application Settings
DEBUG=true
LOG_LEVEL=debug
TLS_CERT=./tls/cert.pem
TLS_KEY=./tls/key.pem

# Docker-specific configuration
MYSQL_ROOT_PASSWORD=dev_root_password_change_me
```

### Manual Setup Environment Variables

For manual setup, use the template to create environment files:

```bash
# Copy template
cp env.template .env.development

# Edit with your settings
# Server Configuration
HOST=localhost
PORT=4000
ENVIROMENT=development

# Database Configuration
DB_DSN=noter_web:pass@tcp(localhost:3306)/noter
TEST_DB_DSN=noter_test_web:test_pass@/test_noter

# Application Settings
DEBUG=true
TLS_CERT=./tls/cert.pem
TLS_KEY=./tls/key.pem
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
