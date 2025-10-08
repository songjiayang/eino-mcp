# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an example project demonstrating integration between CloudWeGo eino and Model Context Protocol (MCP) using both go-mcp and mcp-go implementations. The project creates an AI agent that can use MCP tools to execute tasks.

## Development Commands

### Building and Running
```bash
# Build main application
go build -o eino-mcp main.go

# Run main application (requires MCP server to be running)
go run main.go

# Build MCP time tool (using mcp-go)
cd tools/mcp-time
go build -o mcp-time main.go
./mcp-time

# Build MCP time tool (using go-mcp)
cd tools/mcp-time-v2
go build -o mcp-time-v2 main.go
./mcp-time-v2
```

### Environment Setup
```bash
# Copy environment template and configure
cp env.sh.example env.sh
# Edit env.sh with your API credentials and settings
source env.sh
```

### Testing
```bash
# Run all tests
go test ./...

# Test specific packages
go test ./tools/...
go test ./...
```

## Architecture Overview

### Core Components

**Main Application** (`main.go`):
- Initializes MCP client with SSE transport
- Connects to MCP time tool server at localhost:8080
- Creates eino agent with OpenAI model and MCP tools
- Provides interactive chat interface

**MCP Time Tools** (`tools/`):
- `mcp-time/`: Uses mcp-go library (github.com/mark3labs/mcp-go)
- `mcp-time-v2/`: Uses go-mcp library (github.com/ThinkInAIXYZ/go-mcp)
- Both provide a "current time" tool with timezone support

### Key Dependencies

- **CloudWeGo eino**: Core AI agent framework
- **eino-ext**: Extensions for OpenAI models and MCP tools
- **MCP Libraries**: Two implementations for MCP protocol
- **OpenAI Client**: For LLM integration

### Design Patterns

- **Client-Server Architecture**: Main app connects to MCP tools via SSE
- **Tool Integration**: MCP tools registered as eino agent capabilities
- **Interactive CLI**: Simple command-line interface for agent interaction

## Key Configuration

### Environment Variables
- `OPENAI_API_URL`: OpenAI-compatible API endpoint (default: Aliyun DashScope)
- `MODEL_ID`: Model identifier for the LLM
- `OPENAI_API_KEY`: API key for the service

### MCP Server Setup
- Default MCP server runs on localhost:8080 with SSE transport
- Supports both stdio and SSE transport modes
- Time tool requires timezone parameter for current time queries

## Important Implementation Details

### MCP Client Initialization
```go
// Initialize MCP client with SSE transport
cli, _ := client.NewSSEMCPClient("http://localhost:8080/sse")
cli.Start(ctx)
defer cli.Close()
```

### Agent Configuration
```go
// Create eino agent with MCP tools
llm, _ := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{...})
agent, _ := react.NewAgent(ctx, &react.AgentConfig{
    Model:       llm,
    ToolsConfig: compose.ToolsNodeConfig{Tools: tools},
})
```

### Tool Registration
```go
// Register MCP tools with agent
tools, _ := eino_mcp.GetTools(ctx, &eino_mcp.Config{Cli: cli})
```

## File Structure

- `main.go`: Main application entry point
- `tools/mcp-time/`: MCP time tool implementation (mcp-go)
- `tools/mcp-time-v2/`: MCP time tool implementation (go-mcp)
- `env.sh.example`: Environment configuration template
- `vendor/`: Go module dependencies

## Development Notes

- The project uses Go 1.23.4
- Vendor directory contains all dependencies
- Chinese language interface for user interaction
- OpenTelemetry tracing available through eino-ext components
- Error handling throughout with context cancellation support