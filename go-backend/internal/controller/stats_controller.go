package controller

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type StatsController struct {
	requestLogService *service.RequestLogService
	userService       *service.UserService
}

func NewStatsController(requestLogService *service.RequestLogService, userService *service.UserService) *StatsController {
	return &StatsController{
		requestLogService: requestLogService,
		userService:       userService,
	}
}

func (c *StatsController) GetMyTokenStats(ctx *gin.Context) {
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	totalTokens, err := c.requestLogService.CountUserTokens(loginUser.ID)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	common.Success(ctx, vo.TokenStatsVO{TotalTokens: totalTokens})
}

func (c *StatsController) GetMyLogs(ctx *gin.Context) {
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	limit := 100
	if raw := ctx.Query("limit"); raw != "" {
		if parsed, parseErr := strconv.Atoi(raw); parseErr == nil && parsed > 0 {
			limit = parsed
		}
	}
	logs, err := c.requestLogService.ListUserLogs(loginUser.ID, limit)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	common.Success(ctx, logs)
}

func (c *StatsController) handleError(ctx *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		common.Error(ctx, bizErr.Code, bizErr.Message)
		return
	}
	log.Printf("stats controller error: %v", err)
	common.Error(ctx, errno.SystemError.Code, "系统错误")
}
