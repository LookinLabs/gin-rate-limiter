package corev2

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	ipLimiterInstances = make(map[string]*IPLimiter)
	ipLimiterMutex     sync.Mutex
)

func RequireRateLimiter(rateLimiters ...*RateLimiter) func(*gin.Context) {
	return func(ctx *gin.Context) {
		for _, rateLimiter := range rateLimiters {
			instance, err := getRateLimiterInstance(rateLimiter.RateLimiterType, rateLimiter.Key, rateLimiter.Option)
			if err != nil {
				StatusInternalServerError(ctx, err)
				return
			}

			item, err := instance.GetItem(ctx)
			if err != nil {
				StatusInternalServerError(ctx, err)
				return
			}

			if !item.Limiter.Allow() {
				StatusTooManyRequests(ctx, fmt.Errorf("too many requests"))
				return
			}
		}
	}
}

func getRateLimiterInstance(rateLimiterType RateLimiterType, key string, option RateLimiterOption) (IRateLimiter, error) {
	ipLimiterMutex.Lock()
	defer ipLimiterMutex.Unlock()

	switch rateLimiterType {
	case IPRateLimiter:
		if instance, exists := ipLimiterInstances[key]; exists {
			return instance, nil
		}
		instance := newIPLimiter(key, option)
		ipLimiterInstances[key] = instance
		return instance, nil
	default:
		return nil, fmt.Errorf("rateLimiterType %v is not supported", rateLimiterType)
	}
}
