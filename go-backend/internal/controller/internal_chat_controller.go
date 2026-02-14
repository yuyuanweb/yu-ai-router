package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type InternalChatController struct {
	chatService *service.ChatService
	userService *service.UserService
}

func NewInternalChatController(chatService *service.ChatService, userService *service.UserService) *InternalChatController {
	return &InternalChatController{
		chatService: chatService,
		userService: userService,
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

	apiKeyID := int64(0)
	clientIP := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	if request.Stream != nil && *request.Stream {
		c.stream(ctx, request, loginUser.ID, apiKeyID, clientIP, userAgent)
		return
	}
	response, err := c.chatService.Chat(request, loginUser.ID, apiKeyID, clientIP, userAgent)
	if err != nil {
		writeServiceError(ctx, err)
		return
	}
	common.Success(ctx, response)
}

func (c *InternalChatController) ChatCompletionsWithFile(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("internal chat upload missing file: path=%s err=%v", ctx.Request.URL.Path, err)
		common.Error(ctx, errno.ParamsError.Code, "file 不能为空")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("internal chat upload open file failed: path=%s err=%v", ctx.Request.URL.Path, err)
		common.Error(ctx, errno.OperationError.Code, "读取上传文件失败")
		return
	}
	defer file.Close()

	fileBytes, readErr := io.ReadAll(file)
	if readErr != nil {
		log.Printf("internal chat upload read file failed: path=%s err=%v", ctx.Request.URL.Path, readErr)
		common.Error(ctx, errno.OperationError.Code, "读取上传文件失败")
		return
	}
	if len(fileBytes) == 0 {
		common.Error(ctx, errno.ParamsError.Code, "上传文件不能为空")
		return
	}

	messagesJSON := strings.TrimSpace(ctx.PostForm("messages"))
	if messagesJSON == "" {
		common.Error(ctx, errno.ParamsError.Code, "messages 不能为空")
		return
	}
	var messages []dto.ChatMessage
	if err = json.Unmarshal([]byte(messagesJSON), &messages); err != nil {
		log.Printf("internal chat upload parse messages failed: raw=%s err=%v", messagesJSON, err)
		common.Error(ctx, errno.ParamsError.Code, "messages 格式错误")
		return
	}
	if len(messages) == 0 {
		common.Error(ctx, errno.ParamsError.Code, "messages 不能为空")
		return
	}

	stream := false
	if rawStream := strings.TrimSpace(ctx.PostForm("stream")); rawStream != "" {
		parsed, parseErr := strconv.ParseBool(rawStream)
		if parseErr != nil {
			log.Printf("internal chat upload parse stream failed: raw=%s err=%v", rawStream, parseErr)
			common.Error(ctx, errno.ParamsError.Code, "stream 参数错误")
			return
		}
		stream = parsed
	}

	var enableReasoning *bool
	if rawReasoning := strings.TrimSpace(ctx.PostForm("enable_reasoning")); rawReasoning != "" {
		parsed, parseErr := strconv.ParseBool(rawReasoning)
		if parseErr != nil {
			log.Printf("internal chat upload parse enable_reasoning failed: raw=%s err=%v", rawReasoning, parseErr)
			common.Error(ctx, errno.ParamsError.Code, "enable_reasoning 参数错误")
			return
		}
		enableReasoning = &parsed
	}

	model := strings.TrimSpace(ctx.PostForm("model"))
	if model == "" {
		model = "qwen-plus"
	}
	pluginKey := strings.TrimSpace(ctx.PostForm("plugin_key"))
	contentType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if pluginKey == "" {
		if strings.HasPrefix(contentType, "image/") {
			pluginKey = "image_recognition"
		} else if contentType == "application/pdf" {
			pluginKey = "pdf_parser"
		}
	}

	request := dto.ChatRequest{
		Model:           model,
		Messages:        messages,
		Stream:          &stream,
		RoutingStrategy: strings.TrimSpace(ctx.PostForm("routing_strategy")),
		PluginKey:       pluginKey,
		FileBytes:       fileBytes,
		FileType:        contentType,
		EnableReasoning: enableReasoning,
	}

	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		writeServiceError(ctx, err)
		return
	}
	log.Printf("internal chat upload accepted: userId=%d file=%s size=%d type=%s plugin=%s model=%s stream=%t",
		loginUser.ID, fileHeader.Filename, fileHeader.Size, contentType, pluginKey, model, stream)

	apiKeyID := int64(0)
	clientIP := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")
	if stream {
		c.stream(ctx, request, loginUser.ID, apiKeyID, clientIP, userAgent)
		return
	}
	response, err := c.chatService.Chat(request, loginUser.ID, apiKeyID, clientIP, userAgent)
	if err != nil {
		writeServiceError(ctx, err)
		return
	}
	common.Success(ctx, response)
}

func (c *InternalChatController) stream(ctx *gin.Context, request dto.ChatRequest, userID, apiKeyID int64, clientIP, userAgent string) {
	streamChan, errChan := c.chatService.ChatStream(request, userID, apiKeyID, clientIP, userAgent)
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
				return
			}
			payload, marshalErr := json.Marshal(chunk)
			if marshalErr != nil {
				log.Printf("internal chat stream marshal chunk failed: userId=%d apiKeyId=%d err=%v", userID, apiKeyID, marshalErr)
				continue
			}
			_, _ = ctx.Writer.WriteString("data: " + string(payload) + "\n\n")
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
