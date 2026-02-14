package service

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

const (
	maxFallbackRetries = 3
)

type ChatService struct {
	requestLogService  *RequestLogService
	routingService     *RoutingService
	modelInvokeService *ModelInvokeService
	providerService    *ProviderService
}

func NewChatService(
	requestLogService *RequestLogService,
	routingService *RoutingService,
	modelInvokeService *ModelInvokeService,
	providerService *ProviderService,
) *ChatService {
	return &ChatService{
		requestLogService:  requestLogService,
		routingService:     routingService,
		modelInvokeService: modelInvokeService,
		providerService:    providerService,
	}
}

func (s *ChatService) Chat(chatRequest dto.ChatRequest, userID, apiKeyID int64, clientIP, userAgent string) (*dto.ChatResponse, error) {
	start := time.Now()
	requestedModel := chatRequest.Model
	traceID := newTraceID()
	strategyType := s.routingService.DetermineStrategyType(chatRequest.RoutingStrategy, requestedModel)

	selectedModel, fallbackModels, err := s.routingService.SelectModel(strategyType, constant.ModelTypeChat, requestedModel)
	if err != nil {
		log.Printf("chat select model failed: userId=%d apiKeyId=%d strategy=%s err=%v", userID, apiKeyID, strategyType, err)
		return nil, errno.NewWithMessage(errno.SystemError, "选择模型失败")
	}
	if selectedModel == nil {
		return nil, errno.NewWithMessage(errno.ParamsError, "没有可用的模型")
	}

	candidates := append([]entity.Model{*selectedModel}, fallbackModels...)
	if len(candidates) > maxFallbackRetries+1 {
		candidates = candidates[:maxFallbackRetries+1]
	}
	var lastErr error
	for idx, model := range candidates {
		isFallback := idx > 0
		provider, providerErr := s.providerService.GetProviderByID(model.ProviderID)
		if providerErr != nil {
			lastErr = providerErr
			log.Printf("chat get provider failed: userId=%d apiKeyId=%d model=%s providerId=%d err=%v", userID, apiKeyID, model.ModelKey, model.ProviderID, providerErr)
			continue
		}
		body, invokeErr := s.modelInvokeService.Invoke(&model, provider, withModel(chatRequest, model.ModelKey))
		if invokeErr != nil {
			lastErr = invokeErr
			log.Printf("chat invoke failed: userId=%d apiKeyId=%d model=%s strategy=%s isFallback=%t err=%v", userID, apiKeyID, model.ModelKey, strategyType, isFallback, invokeErr)
			s.requestLogService.LogRequestAsync(RequestLogInput{
				TraceID:         traceID,
				UserID:          ptrInt64(userID),
				APIKeyID:        ptrInt64(apiKeyID),
				ModelID:         ptrInt64(model.ID),
				RequestModel:    model.ModelKey,
				RequestType:     "chat",
				Source:          requestSource(apiKeyID),
				Duration:        int(time.Since(start).Milliseconds()),
				Status:          "failed",
				ErrorMessage:    invokeErr.Error(),
				ErrorCode:       "MODEL_ERROR",
				RoutingStrategy: strategyType,
				IsFallback:      isFallback,
				ClientIP:        clientIP,
				UserAgent:       userAgent,
			})
			continue
		}

		var upstreamResp upstreamChatResponse
		if err = json.Unmarshal(body, &upstreamResp); err != nil {
			lastErr = err
			log.Printf("chat unmarshal upstream response failed: userId=%d apiKeyId=%d model=%s body=%s err=%v", userID, apiKeyID, model.ModelKey, trimForLog(string(body)), err)
			s.requestLogService.LogRequestAsync(RequestLogInput{
				TraceID:         traceID,
				UserID:          ptrInt64(userID),
				APIKeyID:        ptrInt64(apiKeyID),
				ModelID:         ptrInt64(model.ID),
				RequestModel:    model.ModelKey,
				RequestType:     "chat",
				Source:          requestSource(apiKeyID),
				Duration:        int(time.Since(start).Milliseconds()),
				Status:          "failed",
				ErrorMessage:    err.Error(),
				ErrorCode:       "MODEL_ERROR",
				RoutingStrategy: strategyType,
				IsFallback:      isFallback,
				ClientIP:        clientIP,
				UserAgent:       userAgent,
			})
			continue
		}
		response := mapChatResponse(upstreamResp, model.ModelKey)
		usage := response.Usage
		s.requestLogService.LogRequestAsync(RequestLogInput{
			TraceID:          traceID,
			UserID:           ptrInt64(userID),
			APIKeyID:         ptrInt64(apiKeyID),
			ModelID:          ptrInt64(model.ID),
			RequestModel:     model.ModelKey,
			ModelName:        model.ModelKey,
			RequestType:      "chat",
			Source:           requestSource(apiKeyID),
			PromptTokens:     usage.PromptTokens,
			CompletionTokens: usage.CompletionTokens,
			TotalTokens:      usage.TotalTokens,
			Duration:         int(time.Since(start).Milliseconds()),
			Status:           "success",
			RoutingStrategy:  strategyType,
			IsFallback:       isFallback,
			ClientIP:         clientIP,
			UserAgent:        userAgent,
		})
		return &response, nil
	}

	errMsg := "调用模型失败"
	if lastErr != nil {
		errMsg = errMsg + ": " + lastErr.Error()
	}
	s.requestLogService.LogRequestAsync(RequestLogInput{
		TraceID:         traceID,
		UserID:          ptrInt64(userID),
		APIKeyID:        ptrInt64(apiKeyID),
		RequestModel:    requestedModel,
		RequestType:     "chat",
		Source:          requestSource(apiKeyID),
		Duration:        int(time.Since(start).Milliseconds()),
		Status:          "failed",
		ErrorMessage:    errMsg,
		ErrorCode:       "SYSTEM_ERROR",
		RoutingStrategy: strategyType,
		IsFallback:      false,
		ClientIP:        clientIP,
		UserAgent:       userAgent,
	})
	return nil, errno.NewWithMessage(errno.SystemError, errMsg)
}

func (s *ChatService) ChatStream(chatRequest dto.ChatRequest, userID, apiKeyID int64, clientIP, userAgent string) (<-chan dto.StreamResponse, <-chan error) {
	streamChan := make(chan dto.StreamResponse, 32)
	errChan := make(chan error, 1)

	go func() {
		defer close(streamChan)
		defer close(errChan)

		start := time.Now()
		requestedModel := chatRequest.Model
		traceID := newTraceID()
		created := time.Now().Unix()
		strategyType := s.routingService.DetermineStrategyType(chatRequest.RoutingStrategy, requestedModel)
		selectedModel, _, err := s.routingService.SelectModel(strategyType, constant.ModelTypeChat, requestedModel)
		if err != nil {
			log.Printf("chat stream select model failed: userId=%d apiKeyId=%d strategy=%s err=%v", userID, apiKeyID, strategyType, err)
			errChan <- errno.NewWithMessage(errno.SystemError, "选择模型失败")
			return
		}
		if selectedModel == nil {
			errChan <- errno.NewWithMessage(errno.ParamsError, "没有可用的模型")
			return
		}

		provider, providerErr := s.providerService.GetProviderByID(selectedModel.ProviderID)
		if providerErr != nil {
			log.Printf("chat stream get provider failed: userId=%d apiKeyId=%d model=%s providerId=%d err=%v", userID, apiKeyID, selectedModel.ModelKey, selectedModel.ProviderID, providerErr)
			errChan <- errno.NewWithMessage(errno.SystemError, "模型提供者不存在")
			return
		}

		resp, err := s.modelInvokeService.InvokeStream(selectedModel, provider, withModel(chatRequest, selectedModel.ModelKey))
		if err != nil {
			log.Printf("chat stream call upstream failed: userId=%d apiKeyId=%d model=%s err=%v", userID, apiKeyID, selectedModel.ModelKey, err)
			s.requestLogService.LogRequestAsync(RequestLogInput{
				TraceID:         traceID,
				UserID:          ptrInt64(userID),
				APIKeyID:        ptrInt64(apiKeyID),
				ModelID:         ptrInt64(selectedModel.ID),
				RequestModel:    selectedModel.ModelKey,
				RequestType:     "chat",
				Source:          requestSource(apiKeyID),
				Duration:        int(time.Since(start).Milliseconds()),
				Status:          "failed",
				ErrorMessage:    err.Error(),
				ErrorCode:       "STREAM_ERROR",
				RoutingStrategy: strategyType,
				IsFallback:      false,
				ClientIP:        clientIP,
				UserAgent:       userAgent,
			})
			errChan <- errno.NewWithMessage(errno.SystemError, "流式调用模型失败: "+err.Error())
			return
		}
		defer resp.Body.Close()

		promptTokens := 0
		completionTokens := 0
		totalTokens := 0
		isFirstChunk := true
		hasOutput := false
		reader := bufio.NewReader(resp.Body)

		for {
			line, readErr := reader.ReadString('\n')
			if readErr != nil {
				if readErr == io.EOF {
					break
				}
				log.Printf("chat stream read failed: userId=%d apiKeyId=%d model=%s err=%v", userID, apiKeyID, selectedModel.ModelKey, readErr)
				s.requestLogService.LogRequestAsync(RequestLogInput{
					TraceID:         traceID,
					UserID:          ptrInt64(userID),
					APIKeyID:        ptrInt64(apiKeyID),
					ModelID:         ptrInt64(selectedModel.ID),
					RequestModel:    selectedModel.ModelKey,
					RequestType:     "chat",
					Source:          requestSource(apiKeyID),
					Duration:        int(time.Since(start).Milliseconds()),
					Status:          "failed",
					ErrorMessage:    readErr.Error(),
					ErrorCode:       "STREAM_ERROR",
					RoutingStrategy: strategyType,
					IsFallback:      false,
					ClientIP:        clientIP,
					UserAgent:       userAgent,
				})
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
				log.Printf("chat stream chunk unmarshal failed: userId=%d apiKeyId=%d model=%s raw=%s err=%v", userID, apiKeyID, selectedModel.ModelKey, trimForLog(rawData), err)
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

			delta := chunk.Choices[0].Delta
			reasoningContent := strings.TrimSpace(delta.ReasoningContent)
			content := delta.Content
			if reasoningContent == "" && content == "" {
				continue
			}
			responseDelta := dto.StreamResponseDelta{
				Content:          content,
				ReasoningContent: reasoningContent,
			}
			if isFirstChunk {
				responseDelta.Role = "assistant"
				isFirstChunk = false
			}
			streamChan <- dto.StreamResponse{
				ID:      traceID,
				Object:  "chat.completion.chunk",
				Created: created,
				Model:   selectedModel.ModelKey,
				Choices: []dto.StreamResponseChoice{
					{
						Index: 0,
						Delta: responseDelta,
					},
				},
			}
			hasOutput = true
		}

		if !hasOutput {
			log.Printf("chat stream finished with empty output: userId=%d apiKeyId=%d model=%s", userID, apiKeyID, selectedModel.ModelKey)
		}
		streamChan <- dto.StreamResponse{
			ID:      traceID,
			Object:  "chat.completion.chunk",
			Created: created,
			Model:   selectedModel.ModelKey,
			Choices: []dto.StreamResponseChoice{
				{
					Index:        0,
					Delta:        dto.StreamResponseDelta{},
					FinishReason: "stop",
				},
			},
		}

		if totalTokens == 0 {
			totalTokens = promptTokens + completionTokens
		}
		s.requestLogService.LogRequestAsync(RequestLogInput{
			TraceID:          traceID,
			UserID:           ptrInt64(userID),
			APIKeyID:         ptrInt64(apiKeyID),
			ModelID:          ptrInt64(selectedModel.ID),
			RequestModel:     selectedModel.ModelKey,
			ModelName:        selectedModel.ModelKey,
			RequestType:      "chat",
			Source:           requestSource(apiKeyID),
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
			Duration:         int(time.Since(start).Milliseconds()),
			Status:           "success",
			RoutingStrategy:  strategyType,
			IsFallback:       false,
			ClientIP:         clientIP,
			UserAgent:        userAgent,
		})
	}()

	return streamChan, errChan
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
	if v <= 0 {
		return nil
	}
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

func withModel(request dto.ChatRequest, modelName string) dto.ChatRequest {
	request.Model = modelName
	return request
}

func newTraceID() string {
	buffer := make([]byte, 12)
	if _, err := rand.Read(buffer); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(buffer)
}

func requestSource(apiKeyID int64) string {
	if apiKeyID > 0 {
		return "api"
	}
	return "web"
}

type upstreamChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
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
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
