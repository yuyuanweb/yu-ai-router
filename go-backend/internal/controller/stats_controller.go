package controller

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type StatsController struct {
	requestLogService *service.RequestLogService
	userService       *service.UserService
	billingService    *service.BillingService
}

func NewStatsController(
	requestLogService *service.RequestLogService,
	userService *service.UserService,
	billingService *service.BillingService,
) *StatsController {
	return &StatsController{
		requestLogService: requestLogService,
		userService:       userService,
		billingService:    billingService,
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

func (c *StatsController) GetMyCostStats(ctx *gin.Context) {
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	totalCost, totalErr := c.billingService.GetUserTotalCost(loginUser.ID)
	if totalErr != nil {
		c.handleError(ctx, totalErr)
		return
	}
	todayCost, todayErr := c.billingService.GetUserTodayCost(loginUser.ID)
	if todayErr != nil {
		c.handleError(ctx, todayErr)
		return
	}
	common.Success(ctx, vo.CostStatsVO{
		TotalCost: totalCost,
		TodayCost: todayCost,
	})
}

func (c *StatsController) GetMySummaryStats(ctx *gin.Context) {
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	totalTokens, totalTokensErr := c.requestLogService.CountUserTokens(loginUser.ID)
	if totalTokensErr != nil {
		c.handleError(ctx, totalTokensErr)
		return
	}
	remainingQuota, quotaErr := c.userService.GetRemainingQuota(loginUser.ID)
	if quotaErr != nil {
		c.handleError(ctx, quotaErr)
		return
	}
	totalCost, totalCostErr := c.billingService.GetUserTotalCost(loginUser.ID)
	if totalCostErr != nil {
		c.handleError(ctx, totalCostErr)
		return
	}
	todayCost, todayCostErr := c.billingService.GetUserTodayCost(loginUser.ID)
	if todayCostErr != nil {
		c.handleError(ctx, todayCostErr)
		return
	}
	totalRequests, reqErr := c.requestLogService.CountUserRequests(loginUser.ID)
	if reqErr != nil {
		c.handleError(ctx, reqErr)
		return
	}
	successRequests, successErr := c.requestLogService.CountUserSuccessRequests(loginUser.ID)
	if successErr != nil {
		c.handleError(ctx, successErr)
		return
	}
	common.Success(ctx, vo.UserSummaryStatsVO{
		TotalTokens:     totalTokens,
		TokenQuota:      loginUser.TokenQuota,
		UsedTokens:      loginUser.UsedTokens,
		RemainingQuota:  remainingQuota,
		TotalCost:       totalCost,
		TodayCost:       todayCost,
		TotalRequests:   totalRequests,
		SuccessRequests: successRequests,
	})
}

func (c *StatsController) GetMyDailyStats(ctx *gin.Context) {
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	end := time.Now()
	start := end.AddDate(0, 0, -6)
	if raw := ctx.Query("startDate"); raw != "" {
		if parsed, parseErr := time.ParseInLocation("2006-01-02", raw, time.Local); parseErr == nil {
			start = parsed
		}
	}
	if raw := ctx.Query("endDate"); raw != "" {
		if parsed, parseErr := time.ParseInLocation("2006-01-02", raw, time.Local); parseErr == nil {
			end = parsed
		}
	}
	stats, serviceErr := c.requestLogService.GetUserDailyStats(loginUser.ID, start, end)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	common.Success(ctx, stats)
}

func (c *StatsController) PageMyHistory(ctx *gin.Context) {
	var request dto.RequestLogQueryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	loginUser, err := c.userService.GetLoginUser(ctx)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	userID := dto.FlexibleInt64(loginUser.ID)
	request.UserID = &userID
	page, serviceErr := c.requestLogService.PageByQuery(request)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	common.Success(ctx, page)
}

func (c *StatsController) GetHistoryDetail(ctx *gin.Context) {
	id, err := parsePositiveID(ctx.Query("id"))
	if err != nil {
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	loginUser, userErr := c.userService.GetLoginUser(ctx)
	if userErr != nil {
		c.handleError(ctx, userErr)
		return
	}
	record, serviceErr := c.requestLogService.GetByID(id)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	if record == nil {
		common.Error(ctx, errno.NotFoundError.Code, errno.NotFoundError.Message)
		return
	}
	if loginUser.UserRole != "admin" {
		if record.UserID == nil || *record.UserID != loginUser.ID {
			common.Error(ctx, errno.NoAuthError.Code, "只能查看自己的调用历史")
			return
		}
	}
	common.Success(ctx, record)
}

func (c *StatsController) PageHistory(ctx *gin.Context) {
	var request dto.RequestLogQueryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		common.Error(ctx, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	page, serviceErr := c.requestLogService.PageByQuery(request)
	if serviceErr != nil {
		c.handleError(ctx, serviceErr)
		return
	}
	common.Success(ctx, page)
}

func (c *StatsController) handleError(ctx *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		common.Error(ctx, bizErr.Code, bizErr.Message)
		return
	}
	log.Printf("stats controller error: %v", err)
	common.Error(ctx, errno.SystemError.Code, "系统错误")
}
