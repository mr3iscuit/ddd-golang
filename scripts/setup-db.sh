#!/bin/bash

# Database setup script for DDD-Golang Todo Application

echo "Setting up database for DDD-Golang Todo Application..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "Error: docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

echo "Starting PostgreSQL service..."
docker-compose up -d postgres

echo "Waiting for PostgreSQL to be ready..."
sleep 10

# Test database connection
echo "Testing database connection..."
docker-compose exec postgres psql -U todo_user -d todo_db -c "SELECT version();"

if [ $? -eq 0 ]; then
    echo "✅ Database is ready!"
    echo "PostgreSQL is running on localhost:5432"
    echo "Database: todo_db"
    echo "User: todo_user"
    echo "Password: todo_password"
else
    echo "❌ Database connection failed. Please check the logs:"
    docker-compose logs postgres
fi

echo ""
echo "To start the full application:"
echo "  make docker-up"
echo ""
echo "To view logs:"
echo "  make docker-logs"
echo ""
echo "To stop the application:"
echo "  make docker-down" 