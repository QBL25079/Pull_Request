# Makefile
APP_NAME=pr-rewiewer-service

build:
	@echo "Building Go application..."
	go build -o $(APP_NAME) cmd/server/main.go

up:
	@echo "Starting Docker services..."
	docker-compose up -d --build
	@echo "Service is available at http://localhost:8080"

down:
	@echo "Stopping Docker services..."
	docker-compose down

migrate-up:
	@echo "Applying database migrations..."

clean:
	@echo "Cleaning up local build artifacts..."
	rm -f $(APP_NAME)