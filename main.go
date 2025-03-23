package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/cloudwego/eino-ext/components/model/openai"
	eino_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func main() {
	// 初始化 mcp client
	ctx := context.Background()
	cli, _ := client.NewSSEMCPClient("http://localhost:8080/sse")
	cli.Start(ctx)
	defer cli.Close()

	// 发送 init request
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "current-time",
		Version: "1.0.0",
	}
	cli.Initialize(ctx, initRequest)
	// 查询 mcp tools
	tools, _ := eino_mcp.GetTools(ctx, &eino_mcp.Config{Cli: cli})

	// mcp tools 与 eino agent 绑定
	llm, _ := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL: os.Getenv("OPENAI_API_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("OPENAI_API_KEY"),
	})
	agent, _ := react.NewAgent(ctx, &react.AgentConfig{
		Model:       llm,
		ToolsConfig: compose.ToolsNodeConfig{Tools: tools},
	})

	run(agent)
}

func run(agent *react.Agent) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("欢迎使用 eino with mcp demo.")
	inputTips := "\n请输入操作: "
	for {
		fmt.Print(inputTips)
		if !scanner.Scan() {
			fmt.Println("读取输入失败，程序退出。")
			return
		}

		input := scanner.Text()

		switch strings.ToLower(input) {
		case "exit", "bye":
			fmt.Println("欢迎再次使用，再见。")
			return
		default:
			output, err := agent.Generate(context.Background(), []*schema.Message{
				{
					Role:    schema.User,
					Content: strings.Replace(input, inputTips, "", 1),
				},
			})
			if err != nil {
				log.Fatalf("run agent with error %v", err)
			}
			fmt.Println(output.Content)
		}
	}
}
