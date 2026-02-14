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

	resp, err := client.ChatText(context.Background(), "你好，请介绍一下自己")
	if err != nil {
		fmt.Println("调用失败:", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Println("响应:", resp.Choices[0].Message.Content)
	}
	fmt.Println("\nToken 统计:")
	fmt.Println("输入:", resp.Usage.PromptTokens)
	fmt.Println("输出:", resp.Usage.CompletionTokens)
	fmt.Println("总计:", resp.Usage.TotalTokens)
}
