package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic recovered: %v\n%s", recovered, string(debug.Stack()))
				common.Error(c, errno.SystemError.Code, "系统错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}
