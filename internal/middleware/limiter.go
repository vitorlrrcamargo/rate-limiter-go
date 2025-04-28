package middleware

import (
	"fmt"
	"net"
	"strings"

	"rate-limiter-go/internal/limiter"

	"github.com/gin-gonic/gin"
)

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "unknown"
	}
	return parsedIP.String()
}

func getToken(c *gin.Context) string {
	authHeader := c.GetHeader("API_KEY")
	if strings.TrimSpace(authHeader) != "" {
		return authHeader
	}
	return ""
}

func RateLimitMiddleware(limiter limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prioridade: token > IP
		identifier := getToken(c)
		if identifier == "" {
			identifier = getClientIP(c)
		}

		fmt.Println("🔍 Identifier:", identifier)

		allowed, retryAfter, err := limiter.Allow(identifier)
		if err != nil || !allowed {
			c.Header("Retry-After", retryAfter.String())
			c.AbortWithStatusJSON(429, gin.H{
				"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			return
		}

		c.Next()
	}
}
