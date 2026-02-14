package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type ChatController struct {
	chatService   *service.ChatService
	apiKeyService *service.ApiKeyService
}

func NewChatController(chatService *service.ChatService, apiKeyService *service.ApiKeyService) *ChatController {
	return &ChatController{
		chatService:   chatService,
		apiKeyService: apiKeyService,
	}
}

func (c *ChatController) ChatCompletions(ctx *gin.Context) {
	var request dto.ChatRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("chat bind request failed: path=%s err=%v", ctx.Request.URL.Path, err)
		writeError(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	if len(request.Messages) == 0 {
		writeError(ctx, errno.ParamsError.Code, "messages 不能为空")
		return
	}
	authorization := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		log.Printf("chat missing authorization header")
		writeError(ctx, errno.NoAuthError.Code, "缺少或无效的 Authorization Header")
		return
	}
	apiKeyValue := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
	if apiKeyValue == "" {
		log.Printf("chat empty bearer token")
		writeError(ctx, errno.NoAuthError.Code, "缺少或无效的 Authorization Header")
		return
	}
	apiKey, err := c.apiKeyService.GetByKeyValue(apiKeyValue)
	if err != nil {
		writeServiceError(ctx, err)
		return
	}
	if apiKey == nil {
		log.Printf("chat api key invalid or inactive")
		writeError(ctx, errno.NoAuthError.Code, "API Key 无效或已失效")
		return
	}
	clientIP := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	if request.Stream != nil && *request.Stream {
		c.stream(ctx, request, apiKey.UserID, apiKey.ID, clientIP, userAgent)
		return
	}
	response, err := c.chatService.Chat(request, apiKey.UserID, apiKey.ID, clientIP, userAgent)
	if err != nil {
		writeServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *ChatController) stream(ctx *gin.Context, request dto.ChatRequest, userID, apiKeyID int64, clientIP, userAgent string) {
	streamChan, errChan := c.chatService.ChatStream(request, userID, apiKeyID, clientIP, userAgent)
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Status(http.StatusOK)

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		log.Printf("chat stream not supported by writer")
		writeError(ctx, errno.SystemError.Code, "流式响应不支持")
		return
	}

	for {
		select {
		case chunk, open := <-streamChan:
			if !open {
				return
			}
			payload, marshalErr := json.Marshal(chunk)
			if marshalErr != nil {
				log.Printf("chat stream marshal chunk failed: userId=%d apiKeyId=%d err=%v", userID, apiKeyID, marshalErr)
				continue
			}
			_, _ = ctx.Writer.WriteString("data: " + string(payload) + "\n\n")
			flusher.Flush()
		case err, open := <-errChan:
			if open && err != nil {
				log.Printf("chat stream error: %v", err)
				return
			}
		case <-ctx.Request.Context().Done():
			return
		}
	}
}
