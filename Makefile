.PHONY: build clean test run docker docker-build docker-run help

# Default target
all: build

# Build the application
build:
	@echo "Building scallop..."
	@go build -o scallop ./cmd/scallop

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f scallop
	@rm -rf tmp

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the application
run:
	@echo "Running scallop..."
	@./scallop --help

# Run the application with a Docker image
scan:
	@if [ -z "$(IMAGE)" ]; then \
		echo "Usage: make scan IMAGE=<image_name>"; \
		echo "Example: make scan IMAGE=nginx:latest"; \
		exit 1; \
	fi
	@echo "Scanning Docker image: $(IMAGE)..."
	@./scallop --image $(IMAGE) $(ARGS)

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t scallop .

# Run Docker container
docker-run:
	@if [ -z "$(IMAGE)" ]; then \
		echo "Usage: make docker-run IMAGE=<image_name>"; \
		echo "Example: make docker-run IMAGE=nginx:latest"; \
		exit 1; \
	fi
	@echo "Running scallop in Docker container..."
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock scallop --image $(IMAGE) $(ARGS)

# Run Docker Compose
docker-compose:
	@echo "Running scallop with Docker Compose..."
	@docker-compose up --build

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  run           - Run the application"
	@echo "  scan          - Scan a Docker image (make scan IMAGE=nginx:latest)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container (make docker-run IMAGE=nginx:latest)"
	@echo "  docker-compose - Run Docker Compose"
	@echo "  help          - Show this help"
