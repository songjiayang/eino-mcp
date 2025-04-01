package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

var (
	mode string
)

func getTransport() (t transport.ServerTransport) {
	flag.StringVar(&mode, "transport", "stdio", "The transport to use, should be \"stdio\" or \"sse\"")
	flag.Parse()

	if mode == "stdio" {
		log.Println("start mcp server with stdio transport")
		t = transport.NewStdioServerTransport()
	} else {
		addr := "localhost:8080"
		log.Printf("start mcp server with sse transport, listen %s", addr)
		t, _ = transport.NewSSEServerTransport(addr)
	}

	return t
}

func main() {
	server, err := server.NewServer(
		getTransport(),
		server.WithServerInfo(protocol.Implementation{
			Name:    "current-time-v2-client",
			Version: "1.0.0",
		}))

	server.RegisterTool(&protocol.Tool{
		Name:        "current time",
		Description: "Get current time with timezone, Asia/Shanghai is default",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"timezone": map[string]string{
					"type":        "string",
					"description": "current time timezone",
				},
			},
			"required": []string{"timezone"},
		},
	}, func(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		timezone, ok := request.Arguments["timezone"].(string)
		if !ok {
			return nil, errors.New("timezone must be a string")
		}

		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return nil, fmt.Errorf("parse timezone with error: %v", err)
		}
		text := fmt.Sprintf(`current time is %s`, time.Now().In(loc))

		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: text,
				},
			},
		}, err
	})

	if err != nil {
		log.Panicf("sse server start with error: %v", err)
	}
	server.Start()
}
