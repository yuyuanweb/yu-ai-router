package controller

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type ImageController struct {
	imageService  *service.ImageGenerationService
	userService   *service.UserService
	apiKeyService *service.ApiKeyService
}

func NewImageController(
	imageService *service.ImageGenerationService,
	userService *service.UserService,
	apiKeyService *service.ApiKeyService,
) *ImageController {
	return &ImageController{
		imageService:  imageService,
		userService:   userService,
		apiKeyService: apiKeyService,
	}
}

func (i *ImageController) GenerateImage(c *gin.Context) {
	var request dto.ImageGenerationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}

	var userID *int64
	var apiKeyID *int64
	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(authorization, "Bearer ") {
		apiKeyValue := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		apiKey, err := i.apiKeyService.GetByKeyValue(apiKeyValue)
		if err != nil || apiKey == nil {
			common.Error(c, errno.NoAuthError.Code, "API Key 无效或已失效")
			return
		}
		if apiKey.Status != "active" {
			common.Error(c, errno.NoAuthError.Code, "API Key 已被禁用")
			return
		}
		userID = &apiKey.UserID
		apiKeyID = &apiKey.ID
	} else {
		loginUser, err := i.userService.GetLoginUser(c)
		if err != nil {
			common.Error(c, errno.NotLoginError.Code, "请先登录或提供有效的 API Key")
			return
		}
		userID = &loginUser.ID
	}

	response, err := i.imageService.GenerateImage(request, userID, apiKeyID, c.ClientIP())
	if err != nil {
		if bizErr, ok := errno.AsBusinessError(err); ok {
			common.Error(c, bizErr.Code, bizErr.Message)
			return
		}
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	c.JSON(200, response)
}

func (i *ImageController) GetMyRecords(c *gin.Context) {
	loginUser, err := i.userService.GetLoginUser(c)
	if err != nil {
		if bizErr, ok := errno.AsBusinessError(err); ok {
			common.Error(c, bizErr.Code, bizErr.Message)
			return
		}
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	pageNum := int64(1)
	pageSize := int64(10)
	if raw := c.Query("pageNum"); raw != "" {
		if parsed, parseErr := strconv.ParseInt(raw, 10, 64); parseErr == nil && parsed > 0 {
			pageNum = parsed
		}
	}
	if raw := c.Query("pageSize"); raw != "" {
		if parsed, parseErr := strconv.ParseInt(raw, 10, 64); parseErr == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	page, serviceErr := i.imageService.ListUserRecords(loginUser.ID, pageNum, pageSize)
	if serviceErr != nil {
		if bizErr, ok := errno.AsBusinessError(serviceErr); ok {
			common.Error(c, bizErr.Code, bizErr.Message)
			return
		}
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	common.Success(c, page)
}
