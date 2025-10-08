# Eino MCP 集成示例

这是一个展示 CloudWeGo eino 与 Model Context Protocol (MCP) 集成的示例项目。该项目创建了一个 AI 智能体，能够使用 MCP 工具来执行任务。

## 项目特性

- 🤖 基于 CloudWeGo eino 框架的 AI 智能体
- 🔗 MCP (Model Context Protocol) 工具集成
- 🌍 支持多种传输协议 (SSE, stdio)
- ⏰ 实时时间查询工具
- 💬 交互式命令行界面
- 🌐 OpenAI 兼容的 API 集成

## 架构概览

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Eino Agent    │    │   MCP Client     │    │   MCP Tools     │
│                 │◄──►│                  │◄──►│                 │
│ - OpenAI Model  │    │ - SSE Transport   │    │ - Time Tool     │
│ - Tool Registry │    │ - Tool Discovery │    │ - Stdio Mode    │
│ - Chat Interface│    │ - Protocol Handle│    │ - SSE Mode      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 快速开始

### 1. 环境准备

确保你已经安装了 Go 1.23.4 或更高版本：

```bash
go version
```

### 2. 配置环境变量

```bash
# 复制环境配置文件
cp env.sh.example env.sh

# 编辑配置文件，填入你的 API 凭证
cat env.sh
```

配置文件内容：
```bash
export OPENAI_API_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
export MODEL_ID="qwen2.5-32b-instruct"
export OPENAI_API_KEY="your-api-key-here"
```

### 3. 启动 MCP 工具服务器

选择以下任一方式启动 MCP 时间工具：

#### 使用 mcp-go 实现
```bash
cd tools/mcp-time
go build -o mcp-time main.go
./mcp-time -transport=sse -server_listen=localhost:8080
```

#### 使用 go-mcp 实现
```bash
cd tools/mcp-time-v2
go build -o mcp-time-v2 main.go
./mcp-time-v2 -transport=sse
```

### 4. 启动主应用

```bash
# 在项目根目录运行
go run main.go
```

## 使用方法

启动应用后，你将看到交互式界面：

```
欢迎使用 eino with mcp demo.

请输入操作: 现在在北京时间是多少？
当前时间是 2024-01-15 14:30:00 +0800 CST

请输入操作:
```

支持的命令：
- 输入任意自然语言指令来查询时间
- 输入 `exit` 或 `bye` 退出程序

## 项目结构

```
eino-mcp/
├── main.go                    # 主应用程序入口
├── go.mod                     # Go 模块定义
├── go.sum                     # 依赖校验和
├── env.sh.example             # 环境变量模板
├── CLAUDE.md                  # Claude Code 开发指南
├── tools/                     # MCP 工具实现
│   ├── mcp-time/              # mcp-go 实现
│   │   ├── main.go           # 基于 mark3labs/mcp-go
│   │   └── ...               # 相关依赖
│   └── mcp-time-v2/          # go-mcp 实现
│       ├── main.go           # 基于 ThinkInAIXYZ/go-mcp
│       └── ...               # 相关依赖
└── vendor/                    # Go 依赖包
```

## 核心组件

### 主应用 (main.go)

```go
// 初始化 MCP 客户端
cli, _ := client.NewSSEMCPClient("http://localhost:8080/sse")
cli.Start(ctx)

// 创建 OpenAI 模型
llm, _ := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
    BaseURL: os.Getenv("OPENAI_API_URL"),
    Model:   os.Getenv("MODEL_ID"),
    APIKey:  os.Getenv("OPENAI_API_KEY"),
    Timeout: 30 * time.Second,
})

// 创建 eino 智能体
agent, _ := react.NewAgent(ctx, &react.AgentConfig{
    Model:       llm,
    ToolsConfig: compose.ToolsNodeConfig{Tools: tools},
})
```

### MCP 时间工具

两个实现版本：

1. **mcp-time** - 使用 `mcp-go` 库
2. **mcp-time-v2** - 使用 `go-mcp` 库

两个工具都提供相同的功能：
- 获取当前时间
- 支持时区参数
- 返回格式化的时间字符串

## 依赖项

主要依赖：
- `github.com/cloudwego/eino` - 核心 AI 框架
- `github.com/cloudwego/eino-ext` - 扩展组件
- `github.com/mark3labs/mcp-go` - MCP 协议实现 (mcp-time)
- `github.com/ThinkInAIXYZ/go-mcp` - MCP 协议实现 (mcp-time-v2)
- `github.com/cloudwego/eino-ext/components/model/openai` - OpenAI 模型集成

## 开发指南

### 构建应用

```bash
# 构建主应用
go build -o eino-mcp main.go

# 构建工具
cd tools/mcp-time && go build -o mcp-time main.go
cd tools/mcp-time-v2 && go build -o mcp-time-v2 main.go
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 测试特定模块
go test ./tools/...
```

### 添加新工具

1. 在 `tools/` 目录创建新的工具实现
2. 实现 MCP 工具协议
3. 在主应用中注册新工具

## 常见问题

### Q: 如何切换到其他 MCP 传输协议？
A: 使用 `-transport` 参数指定 `stdio` 或 `sse`

### Q: 如何自定义 OpenAI 模型配置？
A: 修改 `env.sh` 文件中的相关环境变量

### Q: 如何添加新的 MCP 工具？
A: 参考 `tools/` 目录下的现有实现，创建新的工具文件

## 许可证

本项目遵循相应的开源许可证。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。