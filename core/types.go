package corev2

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPLimiter struct {
	RateLimiter
}

type RateLimiter struct {
	sync.Mutex
	RateLimiterType RateLimiterType
	Key             string
	Option          RateLimiterOption
	Items           map[string]*RateLimiterItem
}

type IRateLimiter interface {
	GetItem(ctx *gin.Context) (*RateLimiterItem, error)
}

type RateLimiterType int

type RateLimiterItem struct {
	Key        string
	Limiter    *rate.Limiter
	LastSeenAt time.Time
}

type RateLimiterOption struct {
	Limit  rate.Limit
	Burst  int
	Window time.Duration
}

type UpdateRateLimitRequest struct {
	IP     string  `json:"ip"`
	Limit  int     `json:"limit"`
	Burst  int     `json:"burst"`
	Window float64 `json:"window"`
}
