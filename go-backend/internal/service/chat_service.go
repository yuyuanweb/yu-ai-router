package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
)

const (
	httpTimeoutSeconds = 60
	chatPath           = "/v1/chat/completions"
)

type ChatService struct {
	cfg               *config.Config
	requestLogService *RequestLogService
	httpClient        *http.Client
}

func NewChatService(cfg *config.Config, requestLogService *RequestLogService) *ChatService {
	return &ChatService{
		cfg:               cfg,
		requestLogService: requestLogService,
		httpClient: &http.Client{
			Timeout: httpTimeoutSeconds * time.Second,
		},
	}
}

func (s *ChatService) Chat(chatRequest dto.ChatRequest, userID, apiKeyID int64) (*dto.ChatResponse, error) {
	start := time.Now()
	modelName := s.ensureModel(chatRequest.Model)

	payload := buildChatPayload(chatRequest, modelName, false)
	body, err := s.callUpstream(payload)
	if err != nil {
		log.Printf("chat call upstream failed: userId=%d apiKeyId=%d model=%s err=%v", userID, apiKeyID, modelName, err)
		s.requestLogService.LogRequestAsync(ptrInt64(userID), ptrInt64(apiKeyID), modelName, 0, 0, 0, int(time.Since(start).Milliseconds()), "failed", err.Error())
		return nil, errno.NewWithMessage(errno.SystemError, "调用模型失败: "+err.Error())
	}

	var upstreamResp upstreamChatResponse
	if err = json.Unmarshal(body, &upstreamResp); err != nil {
		log.Printf("chat unmarshal upstream response failed: userId=%d apiKeyId=%d model=%s body=%s err=%v", userID, apiKeyID, modelName, trimForLog(string(body)), err)
		s.requestLogService.LogRequestAsync(ptrInt64(userID), ptrInt64(apiKeyID), modelName, 0, 0, 0, int(time.Since(start).Milliseconds()), "failed", err.Error())
		return nil, errno.NewWithMessage(errno.SystemError, "调用模型失败: "+err.Error())
	}

	response := mapChatResponse(upstreamResp, modelName)
	usage := response.Usage
	s.requestLogService.LogRequestAsync(
		ptrInt64(userID),
		ptrInt64(apiKeyID),
		modelName,
		usage.PromptTokens,
		usage.CompletionTokens,
		usage.TotalTokens,
		int(time.Since(start).Milliseconds()),
		"success",
		"",
	)
	return &response, nil
}

func (s *ChatService) ChatStream(chatRequest dto.ChatRequest, userID, apiKeyID int64) (<-chan string, <-chan error) {
	streamChan := make(chan string, 32)
	errChan := make(chan error, 1)

	go func() {
		defer close(streamChan)
		defer close(errChan)

		start := time.Now()
		modelName := s.ensureModel(chatRequest.Model)
		payload := buildChatPayload(chatRequest, modelName, true)
		resp, err := s.callUpstreamStream(payload)
		if err != nil {
			log.Printf("chat stream call upstream failed: userId=%d apiKeyId=%d model=%s err=%v", userID, apiKeyID, modelName, err)
			s.requestLogService.LogRequestAsync(ptrInt64(userID), ptrInt64(apiKeyID), modelName, 0, 0, 0, int(time.Since(start).Milliseconds()), "failed", err.Error())
			errChan <- errno.NewWithMessage(errno.SystemError, "流式调用模型失败: "+err.Error())
			return
		}
		defer resp.Body.Close()

		promptTokens := 0
		completionTokens := 0
		totalTokens := 0
		reader := bufio.NewReader(resp.Body)

		for {
			line, readErr := reader.ReadString('\n')
			if readErr != nil {
				if readErr == io.EOF {
					break
				}
				log.Printf("chat stream read failed: userId=%d apiKeyId=%d model=%s err=%v", userID, apiKeyID, modelName, readErr)
				s.requestLogService.LogRequestAsync(ptrInt64(userID), ptrInt64(apiKeyID), modelName, 0, 0, 0, int(time.Since(start).Milliseconds()), "failed", readErr.Error())
				errChan <- errno.NewWithMessage(errno.SystemError, "流式调用模型失败: "+readErr.Error())
				return
			}
			line = strings.TrimSpace(line)
			if line == "" || !strings.HasPrefix(line, "data:") {
				continue
			}
			rawData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if rawData == "[DONE]" {
				break
			}

			var chunk upstreamStreamChunk
			if err = json.Unmarshal([]byte(rawData), &chunk); err != nil {
				log.Printf("chat stream chunk unmarshal failed: userId=%d apiKeyId=%d model=%s raw=%s err=%v", userID, apiKeyID, modelName, trimForLog(rawData), err)
				continue
			}
			if chunk.Usage.PromptTokens > 0 {
				promptTokens = chunk.Usage.PromptTokens
			}
			if chunk.Usage.CompletionTokens > 0 {
				completionTokens = chunk.Usage.CompletionTokens
			}
			if chunk.Usage.TotalTokens > 0 {
				totalTokens = chunk.Usage.TotalTokens
			}
			if len(chunk.Choices) == 0 {
				continue
			}
			content := chunk.Choices[0].Delta.Content
			if content == "" {
				continue
			}
			streamChan <- strings.ReplaceAll(content, "\n", "\\n")
		}

		if totalTokens == 0 {
			totalTokens = promptTokens + completionTokens
		}
		s.requestLogService.LogRequestAsync(
			ptrInt64(userID),
			ptrInt64(apiKeyID),
			modelName,
			promptTokens,
			completionTokens,
			totalTokens,
			int(time.Since(start).Milliseconds()),
			"success",
			"",
		)
	}()

	return streamChan, errChan
}

func (s *ChatService) callUpstream(payload map[string]any) ([]byte, error) {
	if strings.TrimSpace(s.cfg.AIAPIKey) == "" || strings.Contains(s.cfg.AIAPIKey, "YOUR_QWEN_API_KEY") {
		return nil, errno.NewWithMessage(errno.SystemError, "AI_API_KEY 未配置")
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(s.cfg.AIBaseURL, "/")+chatPath, bytes.NewReader(rawPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.AIAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Printf("chat upstream bad status: status=%d body=%s", resp.StatusCode, trimForLog(string(responseBody)))
		return nil, errno.NewWithMessage(errno.SystemError, string(responseBody))
	}
	return responseBody, nil
}

func (s *ChatService) callUpstreamStream(payload map[string]any) (*http.Response, error) {
	if strings.TrimSpace(s.cfg.AIAPIKey) == "" || strings.Contains(s.cfg.AIAPIKey, "YOUR_QWEN_API_KEY") {
		return nil, errno.NewWithMessage(errno.SystemError, "AI_API_KEY 未配置")
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(s.cfg.AIBaseURL, "/")+chatPath, bytes.NewReader(rawPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.AIAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		defer resp.Body.Close()
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("chat stream upstream bad status: status=%d body=%s", resp.StatusCode, trimForLog(string(responseBody)))
		return nil, errno.NewWithMessage(errno.SystemError, string(responseBody))
	}
	return resp, nil
}

func (s *ChatService) ensureModel(model string) string {
	if strings.TrimSpace(model) == "" {
		return s.cfg.AIModel
	}
	return model
}

func buildChatPayload(request dto.ChatRequest, modelName string, stream bool) map[string]any {
	payload := map[string]any{
		"model":    modelName,
		"messages": request.Messages,
		"stream":   stream,
	}
	if request.Temperature != nil {
		payload["temperature"] = *request.Temperature
	}
	if request.MaxTokens != nil {
		payload["max_tokens"] = *request.MaxTokens
	}
	if stream {
		payload["stream_options"] = map[string]any{
			"include_usage": true,
		}
	}
	return payload
}

func mapChatResponse(upstream upstreamChatResponse, defaultModel string) dto.ChatResponse {
	modelName := upstream.Model
	if modelName == "" {
		modelName = defaultModel
	}
	choices := make([]dto.ChatResponseChoice, 0, len(upstream.Choices))
	for _, item := range upstream.Choices {
		choices = append(choices, dto.ChatResponseChoice{
			Index: item.Index,
			Message: dto.ChatMessage{
				Role:    item.Message.Role,
				Content: item.Message.Content,
			},
			FinishReason: item.FinishReason,
		})
	}
	return dto.ChatResponse{
		ID:      upstream.ID,
		Object:  chooseString(upstream.Object, "chat.completion"),
		Created: upstream.Created,
		Model:   modelName,
		Choices: choices,
		Usage: dto.ChatResponseUsage{
			PromptTokens:     upstream.Usage.PromptTokens,
			CompletionTokens: upstream.Usage.CompletionTokens,
			TotalTokens:      upstream.Usage.TotalTokens,
		},
	}
}

func chooseString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func ptrInt64(v int64) *int64 {
	value := v
	return &value
}

func trimForLog(value string) string {
	const maxLogLength = 1200
	if len(value) <= maxLogLength {
		return value
	}
	return value[:maxLogLength] + "...(truncated)"
}

type upstreamChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int `json:"index"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type upstreamStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
