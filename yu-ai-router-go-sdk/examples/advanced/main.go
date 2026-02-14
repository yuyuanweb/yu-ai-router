package main

import (
	"context"
	"fmt"
	"os"
	"time"

	sdk "github.com/yupi/airouter/yu-ai-router-go-sdk"
)

func main() {
	apiKey := "YOUR_API_KEY"

	cfg := sdk.DefaultConfig(apiKey)
	if baseURL := os.Getenv("YU_AI_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}
	cfg.ConnectTimeout = 15 * time.Second
	cfg.ReadTimeout = 60 * time.Second
	cfg.MaxRetries = 5

	client, err := sdk.NewClient(cfg)
	if err != nil {
		fmt.Println("创建客户端失败:", err)
		return
	}

	temperature := 0.7
	request := sdk.ChatRequest{
		Model: "qwen-turbo",
		Messages: []sdk.ChatMessage{
			sdk.SystemMessage("你是一个编程助手"),
			sdk.UserMessage("什么是 Java？"),
			sdk.AssistantMessage("Java 是一种面向对象的编程语言..."),
			sdk.UserMessage("它的主要特点是什么？"),
		},
		Temperature: &temperature,
	}

	resp, err := client.Chat(context.Background(), request)
	if err != nil {
		fmt.Println("调用失败:", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Println("回答:", resp.Choices[0].Message.Content)
	}
}
