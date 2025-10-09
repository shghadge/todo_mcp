.PHONY: build run test clean mcp-server

# Build the main todo server
build-todo:
	go build -o todo-server main.go

# Build the MCP server
build-mcp:
	go build -o todo-mcp-server cmd/mcp/main.go

# Run the main todo server
run-rest: build-todo
	./todo-server

# Run the MCP server (for testing)
run-mcp: build-mcp
	./todo-mcp-server

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf todo-server todo-mcp-server

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build both servers
build: build-todo build-mcp

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the main todo server"
	@echo "  mcp-server - Build the MCP server"
	@echo "  run        - Run the main todo server"
	@echo "  run-mcp    - Run the MCP server"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  build-all  - Build both servers"
	@echo "  help       - Show this help"
