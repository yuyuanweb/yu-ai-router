package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type UserProviderKeyController struct {
	userProviderKeyService *service.UserProviderKeyService
	userService            *service.UserService
}

func NewUserProviderKeyController(
	userProviderKeyService *service.UserProviderKeyService,
	userService *service.UserService,
) *UserProviderKeyController {
	return &UserProviderKeyController{
		userProviderKeyService: userProviderKeyService,
		userService:            userService,
	}
}

func (c *UserProviderKeyController) AddUserProviderKey(ctx *gin.Context) {
	var request dto.UserProviderKeyAddRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("add byok bind failed: %v", err)
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	result, serviceErr := c.userProviderKeyService.AddUserProviderKey(request, loginUser.ID)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	common.Success(ctx, result)
}

func (c *UserProviderKeyController) UpdateUserProviderKey(ctx *gin.Context) {
	var request dto.UserProviderKeyUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("update byok bind failed: %v", err)
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	result, serviceErr := c.userProviderKeyService.UpdateUserProviderKey(request, loginUser.ID)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	common.Success(ctx, result)
}

func (c *UserProviderKeyController) DeleteUserProviderKey(ctx *gin.Context) {
	var request dto.DeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("delete byok bind failed: %v", err)
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	if request.ID == nil || request.ID.Int64() <= 0 {
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	result, serviceErr := c.userProviderKeyService.DeleteUserProviderKey(request.ID.Int64(), loginUser.ID)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	common.Success(ctx, result)
}

func (c *UserProviderKeyController) ListMyProviderKeys(ctx *gin.Context) {
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	list, serviceErr := c.userProviderKeyService.ListUserProviderKeys(loginUser.ID)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	common.Success(ctx, list)
}

func (c *UserProviderKeyController) handleError(ctx *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		common.Error(ctx, bizErr.Code, bizErr.Message)
		return
	}
	common.Error(ctx, errno.SystemError.Code, "系统错误")
}
