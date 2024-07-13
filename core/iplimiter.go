package core

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (ipl *IPLimiter) GetItem(ctx *gin.Context) (*RateLimiterItem, error) {
	ip := ctx.ClientIP()

	ipl.Lock()
	defer ipl.Unlock()

	if item, exists := ipl.Items[ip]; exists {
		return item, nil
	}

	item := ipl.newItem(ip)
	return item, nil
}

func newIPLimiter(key string, option RateLimiterOption) *IPLimiter {
	return &IPLimiter{
		RateLimiter: RateLimiter{
			RateLimiterType: IPRateLimiter,
			Key:             key,
			Option:          option,
			Items:           make(map[string]*RateLimiterItem),
		},
	}
}

func (ipl *IPLimiter) newItem(ip string) *RateLimiterItem {
	item := &RateLimiterItem{
		Key:        ip,
		Limiter:    rate.NewLimiter(ipl.Option.Limit, ipl.Option.Burst),
		LastSeenAt: time.Now(),
	}
	ipl.Items[ip] = item
	return item
}
