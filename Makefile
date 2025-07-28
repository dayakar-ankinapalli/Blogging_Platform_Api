.PHONY: all build run clean test docker-build docker-run help

BINARY_NAME=go-blog-api

all: build

# Build the Go application
build:
	@echo "Building the application..."
	@go build -o ./bin/${BINARY_NAME} ./cmd/api

# Run the Go application
run:
	@echo "Running the application..."
	@go run ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f ./bin/${BINARY_NAME}

# Build the Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t ${BINARY_NAME} .

# Run the application inside a Docker container
docker-run:
	@echo "Running application in Docker..."
	@docker run --rm -p 8080:8080 ${BINARY_NAME}