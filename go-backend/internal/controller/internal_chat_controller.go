package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type InternalChatController struct {
	chatService   *service.ChatService
	apiKeyService *service.ApiKeyService
	userService   *service.UserService
}

func NewInternalChatController(chatService *service.ChatService, apiKeyService *service.ApiKeyService, userService *service.UserService) *InternalChatController {
	return &InternalChatController{
		chatService:   chatService,
		apiKeyService: apiKeyService,
		userService:   userService,
	}
}

func (c *InternalChatController) ChatCompletions(ctx *gin.Context) {
	var request dto.ChatRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("internal chat bind request failed: path=%s err=%v", ctx.Request.URL.Path, err)
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	if len(request.Messages) == 0 {
		common.Error(ctx, errno.ParamsError.Code, "messages 不能为空")
		return
	}
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		writeServiceError(ctx, err)
		return
	}

	apiKeyID, apiKeyIDProvided, err := parseOptionalPositiveID(ctx.Query("apiKeyId"))
	if err != nil {
		log.Printf("internal chat invalid apiKeyId: raw=%s err=%v", ctx.Query("apiKeyId"), err)
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	if apiKeyIDProvided {
		apiKey, getErr := c.apiKeyService.GetByID(apiKeyID)
		if getErr != nil {
			writeServiceError(ctx, getErr)
			return
		}
		if apiKey == nil {
			log.Printf("internal chat api key not found: apiKeyId=%d userId=%d", apiKeyID, loginUser.ID)
			common.Error(ctx, errno.NotFoundError.Code, "API Key 不存在")
			return
		}
		if apiKey.UserID != loginUser.ID {
			log.Printf("internal chat no auth for api key: apiKeyId=%d loginUserId=%d ownerId=%d", apiKey.ID, loginUser.ID, apiKey.UserID)
			common.Error(ctx, errno.NoAuthError.Code, "无权使用该 API Key")
			return
		}
		if apiKey.Status != "active" {
			log.Printf("internal chat api key inactive: apiKeyId=%d status=%s", apiKey.ID, apiKey.Status)
			common.Error(ctx, errno.ParamsError.Code, "API Key 已失效")
			return
		}
	}

	if request.Stream != nil && *request.Stream {
		c.stream(ctx, request, loginUser.ID, apiKeyID)
		return
	}
	response, err := c.chatService.Chat(request, loginUser.ID, apiKeyID)
	if err != nil {
		writeServiceError(ctx, err)
		return
	}
	common.Success(ctx, response)
}

func (c *InternalChatController) stream(ctx *gin.Context, request dto.ChatRequest, userID, apiKeyID int64) {
	streamChan, errChan := c.chatService.ChatStream(request, userID, apiKeyID)
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Status(http.StatusOK)

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		log.Printf("internal chat stream not supported by writer")
		common.Error(ctx, errno.SystemError.Code, "流式响应不支持")
		return
	}

	for {
		select {
		case chunk, open := <-streamChan:
			if !open {
				_, _ = ctx.Writer.WriteString("data: [DONE]\n\n")
				flusher.Flush()
				return
			}
			_, _ = ctx.Writer.WriteString("data: " + chunk + "\n\n")
			flusher.Flush()
		case err, open := <-errChan:
			if open && err != nil {
				log.Printf("internal chat stream error: userId=%d apiKeyId=%d err=%v", userID, apiKeyID, err)
				return
			}
		case <-ctx.Request.Context().Done():
			return
		}
	}
}

func parseOptionalPositiveID(raw string) (int64, bool, error) {
	if raw == "" {
		return 0, false, nil
	}
	id, err := parsePositiveID(raw)
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}
