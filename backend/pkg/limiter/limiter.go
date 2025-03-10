package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
	Clean()
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
	ttl      time.Duration
}

// NewTokenBucketLimiter 创建令牌桶限流器
// r: 每秒产生的令牌数
// b: 令牌桶容量
// ttl: 限流器生存时间
func NewTokenBucketLimiter(r float64, b int, ttl time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        rate.Limit(r),
		b:        b,
		ttl:      ttl,
	}

	// 启动清理过期限流器的goroutine
	go limiter.cleanupLoop()

	return limiter
}

// getLimiter 获取指定key的限流器
func (l *TokenBucketLimiter) getLimiter(key string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.limiters[key]
	l.mu.RUnlock()

	if !exists {
		l.mu.Lock()
		limiter, exists = l.limiters[key]
		if !exists {
			limiter = rate.NewLimiter(l.r, l.b)
			l.limiters[key] = limiter
		}
		l.mu.Unlock()
	}

	return limiter
}

// Allow 判断是否允许请求通过
func (l *TokenBucketLimiter) Allow(key string) bool {
	return l.getLimiter(key).Allow()
}

// Clean 清理所有限流器
func (l *TokenBucketLimiter) Clean() {
	l.mu.Lock()
	l.limiters = make(map[string]*rate.Limiter)
	l.mu.Unlock()
}

// cleanupLoop 定期清理过期的限流器
func (l *TokenBucketLimiter) cleanupLoop() {
	ticker := time.NewTicker(l.ttl)
	for range ticker.C {
		l.Clean()
	}
}

// IPLimiter IP限流器
type IPLimiter struct {
	limiter *TokenBucketLimiter
}

// NewIPLimiter 创建IP限流器
func NewIPLimiter(r float64, b int, ttl time.Duration) *IPLimiter {
	return &IPLimiter{
		limiter: NewTokenBucketLimiter(r, b, ttl),
	}
}

// Allow 判断IP是否允许请求通过
func (l *IPLimiter) Allow(ip string) bool {
	return l.limiter.Allow(ip)
}

// Clean 清理所有IP限流器
func (l *IPLimiter) Clean() {
	l.limiter.Clean()
}

// APILimiter API限流器
type APILimiter struct {
	limiter *TokenBucketLimiter
}

// NewAPILimiter 创建API限流器
func NewAPILimiter(r float64, b int, ttl time.Duration) *APILimiter {
	return &APILimiter{
		limiter: NewTokenBucketLimiter(r, b, ttl),
	}
}

// Allow 判断API是否允许请求通过
func (l *APILimiter) Allow(path string) bool {
	return l.limiter.Allow(path)
}

// Clean 清理所有API限流器
func (l *APILimiter) Clean() {
	l.limiter.Clean()
}
