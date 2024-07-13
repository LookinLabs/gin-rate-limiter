package example

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/khaaleoo/gin-rate-limiter/core"
)

func Example() {
	router := gin.Default()

	// Create an IP rate limiter middleware
	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip_test",
		Option: ratelimiter.RateLimiterOption{
			Limit: 1,
			Burst: 1,
			Len:   1 * time.Second,
		},
	})

	// Apply rate limiter middleware to a route
	router.GET("/limited-route", rateLimiterMiddleware, func(c *gin.Context) {
		c.String(200, "Hello, rate-limited world!")
	})

	if err := router.Run(":8080"); err != nil {
		log.Printf("failed to start the server: %v", err)
	}
}
