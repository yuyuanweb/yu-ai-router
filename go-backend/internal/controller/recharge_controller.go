package controller

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

const (
	minRechargeAmount = 1.0
	maxRechargeAmount = 10000.0
)

type RechargeController struct {
	rechargeService      *service.RechargeService
	stripePaymentService *service.StripePaymentService
	userService          *service.UserService
	cfg                  *config.Config
}

func NewRechargeController(
	rechargeService *service.RechargeService,
	stripePaymentService *service.StripePaymentService,
	userService *service.UserService,
	cfg *config.Config,
) *RechargeController {
	return &RechargeController{
		rechargeService:      rechargeService,
		stripePaymentService: stripePaymentService,
		userService:          userService,
		cfg:                  cfg,
	}
}

func (r *RechargeController) CreateStripeRecharge(c *gin.Context) {
	var request dto.CreateRechargeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	if request.Amount < minRechargeAmount || request.Amount > maxRechargeAmount {
		common.Error(c, errno.ParamsError.Code, "充值金额必须在1-10000元之间")
		return
	}
	loginUser, err := r.userService.GetLoginUser(c)
	if err != nil {
		if bizErr, ok := errno.AsBusinessError(err); ok {
			common.Error(c, bizErr.Code, bizErr.Message)
			return
		}
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	sess, err := r.stripePaymentService.CreateCheckoutSession(
		loginUser.ID,
		request.Amount,
		r.cfg.StripeSuccessURL,
		r.cfg.StripeCancelURL,
	)
	if err != nil {
		if bizErr, ok := errno.AsBusinessError(err); ok {
			common.Error(c, bizErr.Code, bizErr.Message)
			return
		}
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	common.Success(c, dto.CreateRechargeResponse{
		CheckoutURL: sess.URL,
		SessionID:   sess.ID,
	})
}

func (r *RechargeController) StripeSuccess(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	ok, err := r.stripePaymentService.HandlePaymentSuccess(sessionID)
	if err != nil {
		if bizErr, hasBizErr := errno.AsBusinessError(err); hasBizErr {
			common.Error(c, bizErr.Code, bizErr.Message)
			return
		}
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	if !ok {
		common.Error(c, errno.OperationError.Code, "充值处理失败")
		return
	}
	common.Success(c, "充值成功！")
}

func (r *RechargeController) StripeCancel(c *gin.Context) {
	common.Success(c, "您取消了充值")
}

func (r *RechargeController) GetMyRechargeRecords(c *gin.Context) {
	loginUser, err := r.userService.GetLoginUser(c)
	if err != nil {
		if bizErr, hasBizErr := errno.AsBusinessError(err); hasBizErr {
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
	page, serviceErr := r.rechargeService.ListUserRechargeRecords(loginUser.ID, pageNum, pageSize)
	if serviceErr != nil {
		log.Printf("list recharge records failed: userId=%d err=%v", loginUser.ID, serviceErr)
		common.Error(c, errno.SystemError.Code, errno.SystemError.Message)
		return
	}
	common.Success(c, page)
}
