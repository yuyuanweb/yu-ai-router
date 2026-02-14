package main

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/yupi/airouter/yu-ai-router-go-sdk"
)

func main() {
	apiKey := "YOUR_API_KEY"

	cfg := sdk.DefaultConfig(apiKey)
	if baseURL := os.Getenv("YU_AI_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	client, err := sdk.NewClient(cfg)
	if err != nil {
		fmt.Println("创建客户端失败:", err)
		return
	}

	enableReasoning := true
	temperature := 0.7

	request := sdk.ChatRequest{
		Model: "qwen-plus",
		Messages: []sdk.ChatMessage{
			sdk.UserMessage("请详细解释量子纠缠现象，并给出实际应用"),
		},
		EnableReasoning: &enableReasoning,
		Temperature:     &temperature,
	}

	cb := sdk.FuncStreamCallback{
		OnMessageFunc: func(chunk sdk.ChatChunk) {
			if chunk.ReasoningContent != "" {
				fmt.Printf("\n💭 [思考] %s\n", chunk.ReasoningContent)
			}
			if chunk.Content != "" {
				fmt.Print(chunk.Content)
			}
		},
		OnCompleteFunc: func() {
			fmt.Println("\n\n✅ 完成")
		},
		OnErrorFunc: func(err error) {
			fmt.Println("\n❌ 错误:", err)
		},
	}

	if err = client.ChatStream(context.Background(), request, cb); err != nil {
		fmt.Println("流式调用失败:", err)
	}
}
