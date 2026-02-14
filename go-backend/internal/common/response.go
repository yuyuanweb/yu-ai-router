package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, BaseResponse{
		Code:    0,
		Data:    data,
		Message: "ok",
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, BaseResponse{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}
