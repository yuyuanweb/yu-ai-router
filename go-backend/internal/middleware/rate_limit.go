package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type RateLimitType string

const (
	RateLimitTypeAPIKey RateLimitType = "api_key"
	RateLimitTypeIP     RateLimitType = "ip"
)

func RateLimit(rateLimitService *service.RateLimitService, limitType RateLimitType, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, ok := buildRateLimitKey(c, limitType)
		if !ok {
			c.Next()
			return
		}
		allowed, err := rateLimitService.TryAcquire(key, limit, window)
		if err != nil {
			log.Printf("rate limit check failed: key=%s err=%v", key, err)
			c.Next()
			return
		}
		if !allowed {
			log.Printf("rate limit exceeded: key=%s limit=%d window=%s", key, limit, window.String())
			common.Error(c, errno.TooManyRequest.Code, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

func buildRateLimitKey(c *gin.Context, limitType RateLimitType) (string, bool) {
	switch limitType {
	case RateLimitTypeAPIKey:
		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			clientIP := strings.TrimSpace(c.ClientIP())
			if clientIP == "" {
				return "", false
			}
			return service.BuildIPRateLimitKey(clientIP), true
		}
		apiKey := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		if apiKey == "" {
			clientIP := strings.TrimSpace(c.ClientIP())
			if clientIP == "" {
				return "", false
			}
			return service.BuildIPRateLimitKey(clientIP), true
		}
		return service.BuildAPIKeyRateLimitKey(apiKey), true
	case RateLimitTypeIP:
		clientIP := strings.TrimSpace(c.ClientIP())
		if clientIP == "" {
			return "", false
		}
		return service.BuildIPRateLimitKey(clientIP), true
	default:
		return "", false
	}
}
