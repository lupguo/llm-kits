package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// 初始化 OpenAI 客户端
var client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

// MCP Prompt 示例
func buildMCPPrompt(userInput string) string {
	return fmt.Sprintf(`You are a multi-capability AI Agent.
You have access to the following tools:
- webSearch(query): for recent information
- runCode(code): for calculations
- queryDataStore(query): for structured internal data

Decide which tool(s) to call for the user query, and respond with JSON:
{
  "intent": "webSearch | runCode | queryDataStore",
  "arguments": "...",
  "reasoning": "..."
}

User query: %s
`, userInput)
}

// 简化模拟工具函数
func runCode(input string) string {
	return "🧮 模拟代码运行结果：y = x^2 图像已生成"
}

func webSearch(query string) string {
	return "🌐 模拟网页搜索：2024 Q1 谷歌收入为 800 亿美元"
}

func queryDataStore(query string) string {
	return "📊 模拟数据查询：Q1 销量为 1200 万台"
}

// 调用 OpenAI 并解析计划
func analyzeIntent(userInput string) (string, string, string) {
	prompt := buildMCPPrompt(userInput)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4, // 或 GPT3.5
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "你是一个任务规划AI"},
			{Role: "user", Content: prompt},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	output := resp.Choices[0].Message.Content
	fmt.Println("🧠 MCP计划输出:\n", output)

	// 简单解析（可使用JSON解析替代）
	var intent, args string
	if strings.Contains(output, "webSearch") {
		intent = "webSearch"
	} else if strings.Contains(output, "runCode") {
		intent = "runCode"
	} else if strings.Contains(output, "queryDataStore") {
		intent = "queryDataStore"
	}

	idx := strings.Index(output, "arguments")
	if idx > 0 {
		args = output[idx:]
	}

	return intent, args, output
}

// 入口主函数
func main() {
	var input string
	fmt.Println("请输入任务：")
	fmt.Scanln(&input)

	intent, args, reasoning := analyzeIntent(input)

	var result string
	switch intent {
	case "webSearch":
		result = webSearch(args)
	case "runCode":
		result = runCode(args)
	case "queryDataStore":
		result = queryDataStore(args)
	default:
		result = "无法识别意图"
	}

	fmt.Println("🔎 Agent 分析理由:\n", reasoning)
	fmt.Println("✅ 执行结果:\n", result)
}
