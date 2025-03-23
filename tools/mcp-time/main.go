package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	transport    string
	serverlisten string
)

var (
	serverName    = "github.com/songjiayang/current-time"
	serverVersion = "1.0.0"
)

func init() {
	flag.StringVar(&transport, "transport", "stdio", "The transport to use, should be \"stdio\" or \"sse\"")
	flag.StringVar(&serverlisten, "server_listen", "localhost:8080", "The sse server listen address")
	flag.Parse()
}

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		serverName,
		serverVersion,
	)
	// Add tool
	tool := mcp.NewTool("current time",
		mcp.WithDescription("Get current time with timezone, Asia/Shanghai is default"),
		mcp.WithString("timezone",
			mcp.Required(),
			mcp.Description("current time timezone"),
		),
	)
	// Add tool handler
	s.AddTool(tool, currentTimeHandler)

	// Only check for "sse" since stdio is the default
	if transport == "sse" {
		serverUrl := "http://" + serverlisten
		sseServer := server.NewSSEServer(s, server.WithBaseURL(serverUrl))
		log.Printf("SSE server listening on %s", serverlisten)
		if err := sseServer.Start(serverlisten); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

func currentTimeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	timezone, ok := request.Params.Arguments["timezone"].(string)
	if !ok {
		return nil, errors.New("timezone must be a string")
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("parse timezone with error: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf(`current time is %s`, time.Now().In(loc))), nil
}
