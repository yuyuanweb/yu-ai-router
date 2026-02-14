package yuairoutersdk

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const chatCompletionsPath = "/v1/chat/completions"

type httpClient struct {
	cfg    Config
	client *http.Client
}

func newHTTPClient(cfg Config) *httpClient {
	if cfg.HTTPClient != nil {
		return &httpClient{cfg: cfg, client: cfg.HTTPClient}
	}
	return &httpClient{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.ConnectTimeout + cfg.ReadTimeout,
		},
	}
}

func (c *httpClient) chat(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
	stream := false
	request.Stream = &stream

	var lastErr error
	maxAttempt := c.cfg.MaxRetries + 1
	for attempt := 1; attempt <= maxAttempt; attempt++ {
		response, err := c.doChatRequest(ctx, request)
		if err == nil {
			return response, nil
		}
		lastErr = err

		var authErr *AuthError
		var limitErr *RateLimitError
		if errors.As(err, &authErr) || errors.As(err, &limitErr) {
			return nil, err
		}
		if attempt < maxAttempt && c.cfg.RetryDelay > 0 {
			time.Sleep(c.cfg.RetryDelay * time.Duration(attempt))
		}
	}
	return nil, lastErr
}

func (c *httpClient) doChatRequest(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseURL+chatCompletionsPath, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, mapHTTPError(resp.StatusCode, string(respBody))
	}
	if apiErr := parseAPIError(respBody); apiErr != nil {
		return nil, apiErr
	}

	var chatResp ChatResponse
	if err = json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("parse chat response failed: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("invalid chat response: %s", string(respBody))
	}
	return &chatResp, nil
}

func (c *httpClient) chatStream(ctx context.Context, request ChatRequest, callback StreamCallback) error {
	stream := true
	request.Stream = &stream

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseURL+chatCompletionsPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		callback.OnError(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		err = mapHTTPError(resp.StatusCode, string(respBody))
		callback.OnError(err)
		return err
	}

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, readErr := reader.ReadString('\n')
		if readErr != nil {
			if readErr == io.EOF {
				callback.OnComplete()
				return nil
			}
			callback.OnError(readErr)
			return readErr
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}

		if data == "[DONE]" {
			callback.OnComplete()
			return nil
		}

		if apiErr := parseAPIError([]byte(data)); apiErr != nil {
			callback.OnError(apiErr)
			return apiErr
		}

		var streamResp StreamResponse
		if err = json.Unmarshal([]byte(data), &streamResp); err != nil {
			// 忽略非 JSON 行，保持与 Java SDK 一致的容错策略
			continue
		}
		if len(streamResp.Choices) == 0 {
			continue
		}
		choice := streamResp.Choices[0]
		if choice.FinishReason == "stop" {
			callback.OnComplete()
			return nil
		}
		chunk := ChatChunk{
			Content:          choice.Delta.Content,
			ReasoningContent: choice.Delta.ReasoningContent,
			Model:            streamResp.Model,
			Done:             false,
		}
		if chunk.Content == "" && chunk.ReasoningContent == "" {
			continue
		}
		callback.OnMessage(chunk)
	}
}

type apiErrorEnvelope struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

func parseAPIError(raw []byte) error {
	var envelope apiErrorEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil
	}
	// 避免误判普通 chat 响应，必须同时具备 code 和 message 字段语义
	if envelope.Code == 0 {
		return nil
	}
	message := strings.TrimSpace(envelope.Message)
	if message == "" {
		message = "request failed"
	}
	switch envelope.Code {
	case 40100, 40101:
		return &AuthError{Message: message}
	case 42900:
		return &RateLimitError{Message: message}
	default:
		return &APIError{Code: envelope.Code, Message: message}
	}
}

func mapHTTPError(statusCode int, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		message = fmt.Sprintf("request failed with status %d", statusCode)
	}
	switch statusCode {
	case http.StatusUnauthorized:
		return &AuthError{Message: message}
	case http.StatusTooManyRequests:
		return &RateLimitError{Message: message}
	default:
		return &APIError{Code: statusCode, Message: message}
	}
}
