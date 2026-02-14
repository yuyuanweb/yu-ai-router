package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/service"
)

func RequireLogin(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := userService.GetLoginUser(c); err != nil {
			if bizErr, ok := errno.AsBusinessError(err); ok {
				common.Error(c, bizErr.Code, bizErr.Message)
			} else {
				common.Error(c, errno.SystemError.Code, "系统错误")
			}
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAdmin(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUser, err := userService.GetLoginUser(c)
		if err != nil {
			if bizErr, ok := errno.AsBusinessError(err); ok {
				common.Error(c, bizErr.Code, bizErr.Message)
			} else {
				common.Error(c, errno.SystemError.Code, "系统错误")
			}
			c.Abort()
			return
		}
		if loginUser.UserRole != constant.AdminRole {
			common.Error(c, errno.NoAuthError.Code, errno.NoAuthError.Message)
			c.Abort()
			return
		}
		c.Next()
	}
}
