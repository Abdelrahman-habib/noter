#!/bin/bash

# Docker management script for Snippetbox

set -e

COMPOSE_FILE="docker-compose.yml"
DEV_COMPOSE_FILE="docker-compose.dev.yml"

show_help() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  up          Start all services in production mode"
    echo "  up-dev      Start all services in development mode"
    echo "  down        Stop all services"
    echo "  build       Build the application image"
    echo "  rebuild     Rebuild and restart services"
    echo "  logs        Show logs from all services"
    echo "  logs-web    Show logs from web service only"
    echo "  logs-db     Show logs from database service only"
    echo "  shell       Open shell in web container"
    echo "  db-shell    Open MySQL shell"
    echo "  migrate     Run database migrations manually"
    echo "  test        Run tests in development container"
    echo "  lint        Run linting (go vet)"
    echo "  audit       Run full audit (vet + staticcheck + govulncheck)"
    echo "  tools       Start tools container for development"
    echo "  clean       Remove all containers, volumes, and images"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 up-dev           # Start in development mode"
    echo "  $0 test             # Run tests"
    echo "  $0 audit            # Run security audit"
    echo "  $0 logs-web         # View web application logs"
    echo "  $0 db-shell         # Connect to MySQL database"
}

case "$1" in
    "up")
        echo "Starting Snippetbox in production mode..."
        docker-compose -f $COMPOSE_FILE up -d
        echo "Services started! Access the app at https://localhost:4444"
        ;;
    "up-dev")
        echo "Starting Snippetbox in development mode..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE up -d
        echo "Development services started! Access the app at https://localhost:4444"
        ;;
    "down")
        echo "Stopping all services..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE down
        ;;
    "build")
        echo "Building Snippetbox image..."
        docker-compose -f $COMPOSE_FILE build
        ;;
    "rebuild")
        echo "Rebuilding and restarting services..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE down
        docker-compose -f $COMPOSE_FILE build
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE up -d
        ;;
    "logs")
        docker-compose -f $COMPOSE_FILE logs -f
        ;;
    "logs-web")
        docker-compose -f $COMPOSE_FILE logs -f web
        ;;
    "logs-db")
        docker-compose -f $COMPOSE_FILE logs -f mysql
        ;;
    "shell")
        docker-compose -f $COMPOSE_FILE exec web sh
        ;;
    "db-shell")
        docker-compose -f $COMPOSE_FILE exec mysql mysql -u snippetbox -p snippetbox
        ;;
    "migrate")
        echo "Running database migrations..."
        docker-compose -f $COMPOSE_FILE run --rm migrate
        ;;
    "test")
        echo "Running tests..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE --profile tools run --rm tools go test ./...
        ;;
    "lint")
        echo "Running linting..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE --profile tools run --rm tools go vet ./...
        ;;
    "audit")
        echo "Running full audit..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE --profile tools run --rm tools sh -c "
            echo 'Running go vet...' &&
            go vet ./... &&
            echo 'Running staticcheck...' &&
            staticcheck ./... &&
            echo 'Running govulncheck...' &&
            govulncheck ./...
        "
        ;;
    "tools")
        echo "Starting tools container..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE --profile tools up -d tools
        echo "Tools container started. Use 'docker exec -it snippetbox-tools sh' to access it."
        ;;
    "clean")
        echo "Cleaning up all Docker resources..."
        docker-compose -f $COMPOSE_FILE -f $DEV_COMPOSE_FILE down -v --rmi all
        docker system prune -f
        ;;
    "help"|"")
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac