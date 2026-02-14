package yuairoutersdk

import (
	"context"
	"errors"
)

type StreamCallback interface {
	OnMessage(chunk ChatChunk)
	OnComplete()
	OnError(err error)
}

type Client struct {
	cfg        Config
	httpClient *httpClient
}

func NewClient(cfg Config) (*Client, error) {
	if err := cfg.normalize(); err != nil {
		return nil, err
	}
	return &Client{
		cfg:        cfg,
		httpClient: newHTTPClient(cfg),
	}, nil
}

func (c *Client) Chat(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
	if len(request.Messages) == 0 {
		return nil, errors.New("messages cannot be empty")
	}
	return c.httpClient.chat(ctx, request)
}

func (c *Client) ChatText(ctx context.Context, message string) (*ChatResponse, error) {
	return c.Chat(ctx, SimpleRequest(message))
}

func (c *Client) ChatWithModel(ctx context.Context, model, message string) (*ChatResponse, error) {
	return c.Chat(ctx, RequestWithModel(model, message))
}

func (c *Client) ChatStream(ctx context.Context, request ChatRequest, callback StreamCallback) error {
	if callback == nil {
		return errors.New("callback cannot be nil")
	}
	if len(request.Messages) == 0 {
		return errors.New("messages cannot be empty")
	}
	return c.httpClient.chatStream(ctx, request, callback)
}

func (c *Client) ChatStreamText(ctx context.Context, message string, callback StreamCallback) error {
	return c.ChatStream(ctx, SimpleRequest(message), callback)
}

func (c *Client) ChatStreamWithModel(ctx context.Context, model, message string, callback StreamCallback) error {
	return c.ChatStream(ctx, RequestWithModel(model, message), callback)
}
