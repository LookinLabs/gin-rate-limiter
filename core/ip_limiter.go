package core

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (ipl *IPLimiter) GetConsumerIP(ctx *gin.Context) (*RateLimiterItem, error) {
	consumer := ctx.ClientIP()

	// Use mutex to ensure thread-safe access to ipl.Items map.
	ipl.Lock()
	defer ipl.Unlock()

	if ip, exists := ipl.Items[consumer]; exists {
		return ip, nil
	}

	ip := ipl.newConsumerIP(consumer)
	return ip, nil
}

func newIPRateLimiter(name string, option RateLimiterOption) *IPLimiter {
	return &IPLimiter{
		RateLimiter: RateLimiter{
			RateLimiterType: IPRateLimiter,
			Name:            name,
			Option:          option,
			Items:           make(map[string]*RateLimiterItem),
		},
	}
}

func (ipl *IPLimiter) newConsumerIP(ip string) *RateLimiterItem {
	item := &RateLimiterItem{
		Name:       ip,
		Limiter:    rate.NewLimiter(ipl.Option.Limit, ipl.Option.Burst),
		LastSeenAt: time.Now(),
	}

	ipl.Items[ip] = item
	return item
}
