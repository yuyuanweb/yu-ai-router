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

type BalanceController struct {
	balanceService       *service.BalanceService
	billingRecordService *service.BillingRecordService
	userService          *service.UserService
}

func NewBalanceController(
	balanceService *service.BalanceService,
	billingRecordService *service.BillingRecordService,
	userService *service.UserService,
) *BalanceController {
	return &BalanceController{
		balanceService:       balanceService,
		billingRecordService: billingRecordService,
		userService:          userService,
	}
}

func (b *BalanceController) GetMyBalance(c *gin.Context) {
	loginUser, err := b.userService.GetLoginUser(c)
	if err != nil {
		if bizErr, ok := errno.AsBusinessError(err); ok {
			common.Error(c, bizErr.Code, bizErr.Message)
			return
		}
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	balance, err := b.balanceService.GetUserBalance(loginUser.ID)
	if err != nil {
		b.handleError(c, err)
		return
	}
	totalSpending, err := b.billingRecordService.GetUserTotalSpending(loginUser.ID)
	if err != nil {
		b.handleError(c, err)
		return
	}
	totalRecharge, err := b.billingRecordService.GetUserTotalRecharge(loginUser.ID)
	if err != nil {
		b.handleError(c, err)
		return
	}
	common.Success(c, vo.BalanceVO{
		Balance:       balance,
		TotalSpending: totalSpending,
		TotalRecharge: totalRecharge,
	})
}

func (b *BalanceController) GetMyBillingRecords(c *gin.Context) {
	loginUser, err := b.userService.GetLoginUser(c)
	if err != nil {
		b.handleError(c, err)
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
	page, err := b.billingRecordService.ListUserBillingRecords(loginUser.ID, pageNum, pageSize)
	if err != nil {
		b.handleError(c, err)
		return
	}
	common.Success(c, page)
}

func (b *BalanceController) handleError(c *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		common.Error(c, bizErr.Code, bizErr.Message)
		return
	}
	log.Printf("balance controller error: %v", err)
	common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
}
