package middleware

import (
	"net/http"
	"notex/config"
	"notex/pkg/limiter"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	ipLimiter    *limiter.IPLimiter
	apiLimiter   *limiter.APILimiter
	loginLimiter *limiter.IPLimiter
	limiterMu    sync.RWMutex
)

// InitRateLimiters 初始化限流器
func InitRateLimiters(cfg *config.RateLimitConfig) {
	limiterMu.Lock()
	defer limiterMu.Unlock()

	ipLimiter = limiter.NewIPLimiter(cfg.IP.Rate, cfg.IP.Burst, cfg.IP.TTL)
	apiLimiter = limiter.NewAPILimiter(cfg.API.Rate, cfg.API.Burst, cfg.API.TTL)
	loginLimiter = limiter.NewIPLimiter(cfg.Login.Rate, cfg.Login.Burst, cfg.Login.TTL)
}

// UpdateRateLimiters 更新限流器配置
func UpdateRateLimiters(cfg *config.RateLimitConfig) {
	InitRateLimiters(cfg)
}

// IPRateLimit IP限流中间件
func IPRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiterMu.RLock()
		if ipLimiter == nil {
			InitRateLimiters(&config.GetConfig().RateLimit)
		}
		limiter := ipLimiter
		limiterMu.RUnlock()

		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// APIRateLimit API限流中间件
func APIRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiterMu.RLock()
		if apiLimiter == nil {
			InitRateLimiters(&config.GetConfig().RateLimit)
		}
		limiter := apiLimiter
		limiterMu.RUnlock()

		path := c.FullPath()
		if !limiter.Allow(path) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// LoginRateLimit 登录限流中间件
func LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiterMu.RLock()
		if loginLimiter == nil {
			InitRateLimiters(&config.GetConfig().RateLimit)
		}
		limiter := loginLimiter
		limiterMu.RUnlock()

		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many login attempts, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
