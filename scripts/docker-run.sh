#!/bin/bash

# Simple Docker management script for Noter

set -e

show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  dev         Start in development mode"
    echo "  prod        Start in production mode"
    echo ""
    echo "  migrate         Run database migrations (dev)"
    echo "  migrate-dev     Run database migrations (dev)"
    echo "  migrate-prod    Run database migrations (prod)"
    echo "  migrate-up      Run migrations up (dev)"
    echo "  migrate-down    Run migrations down (dev)"
    echo "  migrate-reset   Reset all migrations (dev)"
    echo "  migrate-status  Check migration status (dev)"
    echo ""
    echo "  tools       Start tools container for development"
    echo "  logs        Show logs from all services"
    echo "  logs-web    Show logs from web service only"
    echo "  logs-db     Show logs from database service only"
    echo "  shell       Open shell in web container"
    echo "  db-shell    Open MySQL shell"
    echo "  down        Stop all services"
    echo "  clean       Remove all containers, volumes, and images"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev              # Start in development mode"
    echo "  $0 prod             # Start in production mode"
    echo "  $0 migrate-dev      # Run migrations (dev)"
    echo "  $0 migrate-prod     # Run migrations (prod)"
    echo "  $0 migrate-down     # Rollback migrations (dev)"
    echo "  $0 logs-web         # View web application logs"
    echo "  $0 db-shell         # Connect to MySQL database"
}

case "$1" in
    "dev")
        echo "Starting Noter in development mode..."
        docker-compose --env-file dev.env up -d
        echo "Development services started! Access the app at https://localhost:4444"
        echo "MySQL is available on localhost:3307"
        ;;
    "prod")
        echo "Starting Noter in production mode..."
        docker-compose --env-file prod.env up -d
        echo "Production services started! Access the app at https://localhost:4444"
        ;;
    "migrate")
        echo "Running database migrations (development)..."
        docker-compose --env-file dev.env --profile migrate up migrate
        ;;
    "migrate-dev")
        echo "Running database migrations (development)..."
        docker-compose --env-file dev.env --profile migrate up migrate
        ;;
    "migrate-prod")
        echo "Running database migrations (production)..."
        docker-compose --env-file prod.env --profile migrate up migrate
        ;;
    "migrate-up")
        echo "Running migrations up (development)..."
        docker-compose --env-file dev.env --profile migrate up migrate
        ;;
    "migrate-down")
        echo "Running migrations down (development)..."
        docker-compose --env-file dev.env --profile migrate run --rm migrate goose -dir db/schema/migrations down
        ;;
    "migrate-reset")
        echo "Resetting all migrations (development)..."
        docker-compose --env-file dev.env --profile migrate run --rm migrate goose -dir db/schema/migrations reset
        ;;
    "migrate-status")
        echo "Checking migration status (development)..."
        docker-compose --env-file dev.env --profile migrate run --rm migrate goose -dir db/schema/migrations status
        ;;
    "tools")
        echo "Starting tools container..."
        docker-compose --profile tools up -d tools
        echo "Tools container started. Use 'docker exec -it noter-tools sh' to access it."
        ;;
    "logs")
        docker-compose logs -f
        ;;
    "logs-web")
        docker-compose logs -f web
        ;;
    "logs-db")
        docker-compose logs -f mysql
        ;;
    "shell")
        docker-compose exec web sh
        ;;
    "db-shell")
        echo "Connecting to MySQL database..."
        docker-compose exec mysql mysql -u root -p${MYSQL_ROOT_PASSWORD:-dev_root_password_change_me}
        ;;
    "down")
        echo "Stopping all services..."
        docker-compose down
        ;;
    "clean")
        echo "Cleaning up all Docker resources..."
        docker-compose down -v --rmi all
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