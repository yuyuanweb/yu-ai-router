package controller

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type UserController struct {
	userService       *service.UserService
	requestLogService *service.RequestLogService
	billingService    *service.BillingService
}

func NewUserController(
	userService *service.UserService,
	requestLogService *service.RequestLogService,
	billingService *service.BillingService,
) *UserController {
	return &UserController{
		userService:       userService,
		requestLogService: requestLogService,
		billingService:    billingService,
	}
}

func (u *UserController) UserRegister(c *gin.Context) {
	var request dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		u.handleBindError(c, "user register", err)
		return
	}
	userID, err := u.userService.UserRegister(request.UserAccount, request.UserPassword, request.CheckPassword)
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, strconv.FormatInt(userID, 10))
}

func (u *UserController) UserLogin(c *gin.Context) {
	var request dto.UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		u.handleBindError(c, "user login", err)
		return
	}
	loginUserVO, err := u.userService.UserLogin(request.UserAccount, request.UserPassword, c)
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, loginUserVO)
}

func (u *UserController) GetLoginUser(c *gin.Context) {
	loginUser, err := u.userService.GetLoginUser(c)
	if err != nil {
		u.handleError(c, err)
		return
	}
	loginUserVO := u.userService.GetLoginUserVO(loginUser)
	common.Success(c, loginUserVO)
}

func (u *UserController) UserLogout(c *gin.Context) {
	if err := u.userService.UserLogout(c); err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, true)
}

func (u *UserController) AddUser(c *gin.Context) {
	var request dto.UserAddRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		u.handleBindError(c, "add user", err)
		return
	}
	userID, err := u.userService.CreateUser(request)
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, strconv.FormatInt(userID, 10))
}

func (u *UserController) GetUserByID(c *gin.Context) {
	id, err := parsePositiveID(c.Query("id"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	user, err := u.userService.GetUserByID(id)
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, user)
}

func (u *UserController) GetUserVOByID(c *gin.Context) {
	id, err := parsePositiveID(c.Query("id"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	user, err := u.userService.GetUserByID(id)
	if err != nil {
		u.handleError(c, err)
		return
	}
	userVO := u.userService.GetUserVO(user)
	common.Success(c, userVO)
}

func (u *UserController) DeleteUser(c *gin.Context) {
	var request dto.DeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		u.handleBindError(c, "delete user", err)
		return
	}
	if request.ID == nil || request.ID.Int64() <= 0 {
		log.Printf("delete user params invalid: id=%v", request.ID)
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, err := u.userService.DeleteUser(request.ID.Int64())
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (u *UserController) UpdateUser(c *gin.Context) {
	var request dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		u.handleBindError(c, "update user", err)
		return
	}
	if request.ID == nil {
		log.Printf("update user params invalid: missing id")
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, err := u.userService.UpdateUser(request)
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (u *UserController) ListUserVOByPage(c *gin.Context) {
	var request dto.UserQueryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		u.handleBindError(c, "list user by page", err)
		return
	}
	pageResponse, err := u.userService.ListUserVOByPage(request)
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, pageResponse)
}

func (u *UserController) GetMyQuota(c *gin.Context) {
	loginUser, err := u.userService.GetLoginUser(c)
	if err != nil {
		u.handleError(c, err)
		return
	}
	remainingQuota, err := u.userService.GetRemainingQuota(loginUser.ID)
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, vo.QuotaVO{
		TokenQuota:     loginUser.TokenQuota,
		UsedTokens:     loginUser.UsedTokens,
		RemainingQuota: remainingQuota,
	})
}

func (u *UserController) SetUserQuota(c *gin.Context) {
	var request dto.QuotaUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		u.handleBindError(c, "set user quota", err)
		return
	}
	if request.UserID == nil || request.TokenQuota == nil || request.UserID.Int64() <= 0 {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, err := u.userService.SetUserQuota(request.UserID.Int64(), request.TokenQuota.Int64())
	if err != nil {
		u.handleError(c, err)
		return
	}
	common.Success(c, result)
}

func (u *UserController) ResetUserQuota(c *gin.Context) {
	userID, err := parsePositiveID(c.Query("userId"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, serviceErr := u.userService.ResetUserQuota(userID)
	if serviceErr != nil {
		u.handleError(c, serviceErr)
		return
	}
	common.Success(c, result)
}

func (u *UserController) DisableUser(c *gin.Context) {
	userID, err := parsePositiveID(c.Query("userId"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, serviceErr := u.userService.DisableUser(userID)
	if serviceErr != nil {
		u.handleError(c, serviceErr)
		return
	}
	common.Success(c, result)
}

func (u *UserController) EnableUser(c *gin.Context) {
	userID, err := parsePositiveID(c.Query("userId"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	result, serviceErr := u.userService.EnableUser(userID)
	if serviceErr != nil {
		u.handleError(c, serviceErr)
		return
	}
	common.Success(c, result)
}

func (u *UserController) GetUserAnalysis(c *gin.Context) {
	userID, err := parsePositiveID(c.Query("userId"))
	if err != nil {
		common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
		return
	}
	user, serviceErr := u.userService.GetUserByID(userID)
	if serviceErr != nil {
		u.handleError(c, serviceErr)
		return
	}
	remainingQuota, quotaErr := u.userService.GetRemainingQuota(userID)
	if quotaErr != nil {
		u.handleError(c, quotaErr)
		return
	}
	totalRequests, reqErr := u.requestLogService.CountUserRequests(userID)
	if reqErr != nil {
		u.handleError(c, reqErr)
		return
	}
	successRequests, successErr := u.requestLogService.CountUserSuccessRequests(userID)
	if successErr != nil {
		u.handleError(c, successErr)
		return
	}
	totalTokens, tokensErr := u.requestLogService.CountUserTokens(userID)
	if tokensErr != nil {
		u.handleError(c, tokensErr)
		return
	}
	totalCost, totalCostErr := u.billingService.GetUserTotalCost(userID)
	if totalCostErr != nil {
		u.handleError(c, totalCostErr)
		return
	}
	todayCost, todayCostErr := u.billingService.GetUserTodayCost(userID)
	if todayCostErr != nil {
		u.handleError(c, todayCostErr)
		return
	}
	common.Success(c, vo.UserAnalysisVO{
		UserID:          user.ID,
		UserAccount:     user.UserAccount,
		UserName:        user.UserName,
		UserStatus:      user.UserStatus,
		UserRole:        user.UserRole,
		TokenQuota:      user.TokenQuota,
		UsedTokens:      user.UsedTokens,
		RemainingQuota:  remainingQuota,
		TotalRequests:   totalRequests,
		SuccessRequests: successRequests,
		TotalTokens:     totalTokens,
		TotalCost:       totalCost,
		TodayCost:       todayCost,
	})
}

func (u *UserController) handleError(c *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		log.Printf("business error: method=%s path=%s code=%d message=%s", c.Request.Method, c.Request.URL.Path, bizErr.Code, bizErr.Message)
		common.Error(c, bizErr.Code, bizErr.Message)
		return
	}
	log.Printf("system error: method=%s path=%s err=%v", c.Request.Method, c.Request.URL.Path, err)
	common.Error(c, errno.SystemError.Code, "系统错误")
}

func (u *UserController) handleBindError(c *gin.Context, action string, err error) {
	log.Printf("bind request failed: action=%s method=%s path=%s err=%v", action, c.Request.Method, c.Request.URL.Path, err)
	common.Error(c, errno.ParamsError.Code, errno.ParamsError.Message)
}

func parsePositiveID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errno.New(errno.ParamsError)
	}
	return id, nil
}
