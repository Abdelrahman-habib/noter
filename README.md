# Noter

A secure web application for sharing Notes, Quotes, and any text in general, built with Go and PostgreSQL. Features user authentication, session management, and HTTPS support.

## Features

- 📝 Create and share text notes
- 🔐 User registration and authentication
- 🍪 Secure session management with PostgreSQL store
- 🔒 CSRF protection and bcrypt password hashing
- 🌐 HTTPS/TLS support
- 📱 Responsive web interface
- 🐳 Docker support for easy deployment

## Technology Stack

- **Language**: Go 1.25.1
- **Database**: PostgreSQL 17 (via the `pgx/v5/stdlib` database/sql driver)
- **Web Framework**: Standard library `net/http` with custom routing
- **Template Engine**: Go's `html/template`
- **Session Management**: SCS with PostgreSQL store
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
   cp env.template dev.env
   # In dev.env, change localhost:5433 to postgres:5432 in all database URLs.
   ./scripts/docker-run.sh dev
   ```

3. **Run database migrations** (first time only)

   ```bash
   ./scripts/docker-run.sh migrate
   ```

4. **Access the application**
   - Open https://localhost:4444 in your browser
   - Accept the self-signed certificate warning for development
   - PostgreSQL is available on localhost:5433

### Production Deployment

```bash
./docker-run.sh prod
```

## Docker Setup Details

### Architecture

The Docker setup uses a multi-container architecture:

- **PostgreSQL Container**: Database server with automatic initialization
- **Web Container**: Go application with live reload for development
- **Migration Container**: Runs Goose migrations (on-demand)

### Services

#### PostgreSQL Service (`noter-postgres`)

- **Image**: `postgres:17`
- **Port**: `5433:5432` (mapped to avoid conflicts with a local PostgreSQL)
- **Data**: Persistent volume `pg_data`
- **Initialization**: Automatically creates roles and databases via `db/init/`
- **Health Check**: Ensures PostgreSQL is ready before other services start

#### Web Service (`noter-web`)

- **Build**: Custom Go application image
- **Port**: `4444:4444` (HTTPS)
- **Volumes**:
  - Source code mounted for live reload
  - Go module cache for faster builds
- **Dependencies**: Waits for PostgreSQL to be healthy
- **Environment**: Configurable via environment variables

#### Migration Service (`noter-migrate`)

- **Profile**: `migrate` (runs on-demand)
- **Purpose**: Runs Goose migrations and seed data
- **Dependencies**: Waits for PostgreSQL to be healthy
- **Command**: `goose -dir db/schema/migrations up && goose -dir db/schema/seed -no-versioning up`

### Database Initialization

The PostgreSQL container creates the databases and roles. The migration service
creates tables and loads development seed data:

1. **Database Creation**: `noter` and `noter_test` databases
2. **Role Creation**: Three roles with appropriate permissions:
   - `noter_admin`: Owner of the `noter` database (runs migrations)
   - `noter_web`: Read/write access for the application
   - `noter_test_web`: Owner of the `noter_test` database (for tests)
3. **Migration Execution**: The migration service runs Goose to create tables
4. **Seed Data**: Sample data inserted for development

### Environment Configuration

#### Development Environment (`dev.env`)

```env
# Server Configuration
HOST=localhost
PORT=4000
ENVIROMENT=development

# Database Configuration
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://noter_admin:admin@postgres:5432/noter?sslmode=disable
GOOSE_MIGRATION_DIR=./db/schema/migrations
GOOSE_SEED_DIR=./db/schema/seed
GOOSE_TABLE=goose_migrations
DB_DSN=postgres://noter_web:pass@postgres:5432/noter?sslmode=disable
TEST_DB_DSN=postgres://noter_test_web:test_pass@postgres:5432/noter_test?sslmode=disable

# Application Settings
DEBUG=true
LOG_LEVEL=debug
TLS_CERT=./tls/cert.pem
TLS_KEY=./tls/key.pem

# Docker-specific configuration
POSTGRES_PASSWORD=postgres_dev_change_me
```

Use `postgres:5432` for connections between Docker services and `localhost:5433`
when connecting to the Docker database from the host.

#### Security Considerations

- **PostgreSQL Superuser Password**: Configurable via `POSTGRES_PASSWORD` environment variable
- **Default Password**: `postgres_dev_change_me` (change for anything shared)
- **Database Roles**: Separate roles with minimal required permissions
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
./docker-run.sh db-shell         # Access PostgreSQL shell
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

   - PostgreSQL uses port 5433 to avoid conflicts with a local PostgreSQL
   - Web app uses port 4444 for HTTPS

2. **Database Connection Issues**

   - Ensure PostgreSQL container is healthy before running migrations
   - Check that database roles have correct permissions

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
- PostgreSQL 17 or later
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
   # Create roles and databases (dev-only default passwords):
   psql -U postgres -f db/init/01-create-databases-and-users.sql

   # Copy development configuration
   cp env.template development.env
   # Edit development.env with your database credentials

   # Run migrations
   make migrate-up ENV=development

   # Seed development data (optional)
   make seed-up ENV=development
   ```

   For a local PostgreSQL installation, use its port (usually `5432`) in the
   database URLs. The template uses `5433` for the Docker database.

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
make test             # Run all tests (models tests need a PostgreSQL test database)
make test-short       # Run tests with the database-backed tests skipped
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
./scripts/docker-run.sh logs-db     # View database logs
./scripts/docker-run.sh shell       # Access app container
./scripts/docker-run.sh db-shell    # Access PostgreSQL shell

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
│   └── init/            # PostgreSQL initialization scripts
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

Create configuration files from `env.template` for Make or Docker to load:

- `env.template` - Template for creating environment files
- `dev.env`, `prod.env` - Docker settings
- `development.env`, `production.env`, `test.env` - Make settings, selected with `ENV`
- `.env` - Make fallback when the selected environment file is missing

The app takes its defaults from the environment — `PORT`, `ENVIROMENT`,
`DEBUG`, `TLS_CERT`, `TLS_KEY`, `DB_DSN` — and an explicitly passed flag wins
over them. Goose uses `GOOSE_DRIVER=postgres` and `GOOSE_DBSTRING` for migrations.

### Docker Environment Variables

The Docker setup uses `dev.env` for development configuration:

```bash
# Server Configuration
HOST=localhost
PORT=4000
ENVIROMENT=development

# Database Configuration
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://noter_admin:admin@postgres:5432/noter?sslmode=disable
GOOSE_MIGRATION_DIR=./db/schema/migrations
GOOSE_SEED_DIR=./db/schema/seed
GOOSE_TABLE=goose_migrations
DB_DSN=postgres://noter_web:pass@postgres:5432/noter?sslmode=disable
TEST_DB_DSN=postgres://noter_test_web:test_pass@postgres:5432/noter_test?sslmode=disable

# Application Settings
DEBUG=true
LOG_LEVEL=debug
TLS_CERT=./tls/cert.pem
TLS_KEY=./tls/key.pem

# Docker-specific configuration
POSTGRES_PASSWORD=postgres_dev_change_me
```

### Manual Setup Environment Variables

For manual setup, use the template to create environment files:

```bash
# Copy development configuration
cp env.template development.env

# Edit with your settings
# Server Configuration
HOST=localhost
PORT=4000
ENVIROMENT=development

# Database Configuration (local PostgreSQL default port 5432; use 5433 for the Docker-mapped port)
DB_DSN=postgres://noter_web:pass@localhost:5432/noter?sslmode=disable
TEST_DB_DSN=postgres://noter_test_web:test_pass@localhost:5432/noter_test?sslmode=disable

# Application Settings
DEBUG=true
TLS_CERT=./tls/cert.pem
TLS_KEY=./tls/key.pem
```

## Security Features

- **CSRF Protection**: All forms protected against cross-site request forgery
- **Secure Sessions**: HTTP-only, secure cookies with PostgreSQL storage
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
docker-compose --env-file dev.env up -d postgres
docker-compose --env-file dev.env --profile tools run --rm tools go test ./...
```

### Test Database Setup

The bootstrap script in `db/init/` creates the `noter_test` database and its
owner. Model tests apply and roll back their own migrations and seeds, so
start with an empty test schema. Configure `TEST_DB_DSN` for this database;
the default uses `localhost:5433`.

`make init-test-db` and `make teardown-test-db` are available for manual schema
inspection. Tear down a manually initialized schema before running the tests.

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
