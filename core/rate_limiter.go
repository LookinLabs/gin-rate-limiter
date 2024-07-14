package core

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	ipLimiterMap   = make(map[string]*IPLimiter)
	ipLimiterMutex sync.Mutex
)

func RequireRateLimiter(rateLimiters ...*RateLimiter) func(*gin.Context) {
	return func(ctx *gin.Context) {
		for _, rateLimiter := range rateLimiters {
			instance, err := getRateLimiterInstance(rateLimiter.RateLimiterType, rateLimiter.Name, rateLimiter.Option)
			if err != nil {
				StatusInternalServerError(ctx, err)
				return
			}

			allowedIP, err := instance.GetConsumerIP(ctx)
			if err != nil {
				StatusInternalServerError(ctx, err)
				return
			}

			if !allowedIP.Limiter.Allow() {
				StatusTooManyRequests(ctx, fmt.Errorf("too many requests"))
				return
			}
		}
	}
}

func getRateLimiterInstance(rateLimiterType RateLimiterType, ip string, option RateLimiterOption) (IRateLimiter, error) {
	// Use mutex to ensure thread-safe access to ipl.Items map.
	ipLimiterMutex.Lock()
	defer ipLimiterMutex.Unlock()

	switch rateLimiterType {
	case IPRateLimiter:
		if ipRateLimiter, exists := ipLimiterMap[ip]; exists {
			return ipRateLimiter, nil
		}

		ipRateLimiter := newIPRateLimiter(ip, option)
		ipLimiterMap[ip] = ipRateLimiter

		return ipRateLimiter, nil

	default:
		return nil, fmt.Errorf("rateLimiterType %v is not supported", rateLimiterType)
	}
}
