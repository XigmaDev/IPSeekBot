BINARY_NAME=IPSeekBot

GO_FILES=$(shell find . -name "*.go")

GO_MOD=go.mod
GO_SUM=go.sum

.PHONY: all
all: run

.PHONY: run
run:
	@echo "Running the bot..."
	@go run .


.PHONY: tidy
tidy:
	@go mod tidy -v -x 


.PHONY: build
build:
	@echo "Building the bot binary..."
	@go build -o bin/$(BINARY_NAME)

.PHONY: start
start:
	@echo "Starting the bot binary..."
	@./bin/$(BINARY_NAME)


.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@gofmt -w .

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f bin/$(BINARY_NAME)

.PHONY: deps
deps:
	@echo "Installing project dependencies..."
	@go mod tidy

.PHONY: help
help:
	@echo "Available make commands:"
	@echo "  run    - Run the bot"
	@echo "  build  - Build the bot binary"
	@echo "  start  - Run the bot binary"
	@echo "  lint   - Run golangci-lint"
	@echo "  fmt    - Format code using gofmt"
	@echo "  test   - Run unit tests"
	@echo "  clean  - Remove the built binary"
	@echo "  deps   - Install project dependencies"
	@echo "  help   - Display this help message"
