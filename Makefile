# Variables
ENV_FILE = .env
GO_CMD = CGO_ENABLED=1 GOOS=linux go

# Build the Feedgen Go binary
build:
	@echo "Building Feed Generator Go binary..."
	$(GO_CMD) build -o feedgen cmd/*.go

up:
	@echo "Starting Go Feed Generator..."
	docker compose -f docker-compose.yml up --build -d
