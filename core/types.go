package core

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPLimiter struct {
	RateLimiter
}

type IRateLimiter interface {
	GetConsumerIP(ctx *gin.Context) (*RateLimiterItem, error)
}

type RateLimiter struct {
	sync.Mutex
	RateLimiterType RateLimiterType
	Name            string
	Option          RateLimiterOption
	Items           map[string]*RateLimiterItem
}

type RateLimiterType int

type RateLimiterItem struct {
	Name       string
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
