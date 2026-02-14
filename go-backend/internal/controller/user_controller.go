package controller

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
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
