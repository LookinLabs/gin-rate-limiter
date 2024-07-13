package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/lookinlabs/gin-rate-limiter/core"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	// Create an IP rate limiter middleware
	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(&ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip_test",
		Option: ratelimiter.RateLimiterOption{
			Limit:  1,
			Burst:  500,
			Window: 10 * time.Minute,
		},
	})

	return rateLimiterMiddleware
}

func main() {
	router := gin.Default()

	// Apply the rate limiter middleware
	router.GET("/me", RateLimiterMiddleware(), func(ctx *gin.Context) {
		ratelimiter.StatusOK(ctx, gin.H{"message": "hello world"})
	})

	router.Run(":8080")
}
