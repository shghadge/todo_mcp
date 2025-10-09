# Todo MCP Server

A comprehensive todo application with both REST API and MCP (Model Context Protocol) server implementations for LLM integration.

## Features

### Todo Application
- **CRUD Operations**: Create, read, update, and delete todos
- **Status Management**: Todos can be pending or completed
- **REST API**: Full RESTful API
- **In-Memory Storage**: Thread-safe in-memory storage with concurrent access

### MCP Server
- **Model Context Protocol**: Full MCP server implementation for LLM integration
- **Tools**: 5 interactive tools for todo management
- **Resources**: 3 resources for accessing todo data
- **JSON-RPC Protocol**: Standard MCP communication protocol

## Architecture

```
todo_mcp/
├── cmd/
│   └── mcp/          # MCP server entry point 
├── internal/
│   ├── handlers/     # HTTP handlers for REST API
│   ├── models/       # Data models
│   ├── mcp/          # MCP server implementation
│   └── storage/      # Storage interface and implementations
├── main.go           # REST API server
├── mcp_config.json   # MCP server configuration
└── Makefile          # Build automation
```

## MCP Server Capabilities

### Tools
1. **create_todo** - Create a new todo item
2. **get_todo** - Get a specific todo by ID
3. **get_todos** - Get all todos or filter by status
4. **update_todo** - Update an existing todo
5. **delete_todo** - Delete a todo by ID

### Resources
1. **todo://todos** - All todos
2. **todo://todos/pending** - Pending todos only
3. **todo://todos/completed** - Completed todos only

## Quick Start

### Prerequisites
- Go 1.23.2 or later

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd todo_mcp
```

2. Install dependencies:
```bash
make deps
```

3. Build both servers:
```bash
make build-all
```

### Running the REST API Server

```bash
make run
```

The REST API will be available at `http://localhost:8080`

### Running the MCP Server

```bash
make run-mcp
```

The MCP server communicates via stdio using JSON-RPC protocol.

## Usage

### REST API Endpoints

- `POST /api/v1/todos` - Create a todo
- `GET /api/v1/todos` - Get all todos (optional ?status=pending|completed)
- `GET /api/v1/todos/{id}` - Get a specific todo
- `PUT /api/v1/todos/{id}` - Update a todo
- `DELETE /api/v1/todos/{id}` - Delete a todo

### MCP Server Integration

#### For Claude Desktop
Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "todo-mcp-server": {
      "command": "/path/to/todo_mcp/todo-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

#### For Other MCP Clients
Use the MCP server binary with stdio transport:

```bash
./todo-mcp-server
```

## Testing

### Test the REST API
```bash
# Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Todo","description":"This is a test"}'

# Get all todos
curl http://localhost:8080/api/v1/todos

# Get pending todos
curl http://localhost:8080/api/v1/todos?status=pending
```

### Test the MCP Server
```bash
# Build test client
go build -o test-mcp-client cmd/test-mcp/main.go

# Run test client (requires MCP server to be built first)
make mcp-server
./test-mcp-client
```

## MCP Protocol Details

### Tool Examples

#### Create Todo
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_todo",
    "arguments": {
      "title": "Buy groceries",
      "description": "Pick up milk and bread"
    }
  }
}
```

#### Get Todos
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get_todos",
    "arguments": {
      "status": "pending"
    }
  }
}
```

#### Update Todo
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "update_todo",
    "arguments": {
      "id": 1,
      "status": "completed"
    }
  }
}
```

### Resource Examples

#### List Resources
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "resources/list",
  "params": {}
}
```

#### Read Resource
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "resources/read",
  "params": {
    "uri": "todo://todos/pending"
  }
}
```

## Development

### Project Structure
- `internal/models/` - Data models and request/response types
- `internal/storage/` - Storage interface and implementations
- `internal/handlers/` - HTTP request handlers
- `internal/mcp/` - MCP server implementation
  - `models.go` - MCP protocol types
  - `server.go` - Main MCP server logic
  - `tools.go` - Tool implementations
  - `resources.go` - Resource implementations

### Adding New Features
1. Add new models to `internal/models/`
2. Implement storage methods in `internal/storage/`
3. Add HTTP handlers in `internal/handlers/`
4. For MCP: Add tools in `internal/mcp/tools.go` and resources in `internal/mcp/resources.go`

### Build Targets
- `make build` - Build REST API server
- `make mcp-server` - Build MCP server
- `make build-all` - Build both servers
- `make test` - Run tests
- `make clean` - Clean build artifacts

## License

This project is for learning purposes and demonstrates MCP server implementation.
