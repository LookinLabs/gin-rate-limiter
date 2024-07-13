package tests

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/lookinlabs/gin-rate-limiter/core"
	"github.com/stretchr/testify/assert"
)

var wg sync.WaitGroup

func TestMaximumRequestsInAPeriod(testCase *testing.T) {
	router := SetupRouter()

	// Rate limiter configuration
	WindowCapacity := 5
	WindowLen := 1 * time.Second

	// Fake requests
	successCount := 0
	errorCount := 0
	numRequestsPerRoute := 5
	requests := make([]int, numRequestsPerRoute)

	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(&ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip_test",
		Option: ratelimiter.RateLimiterOption{
			Limit:  1,
			Burst:  WindowCapacity,
			Window: WindowLen,
		},
	})

	router.GET("/me", rateLimiterMiddleware, func(ctx *gin.Context) {
		ratelimiter.StatusOK(ctx, gin.H{"message": "hello world"})
	})

	timeStart := time.Now()

	for i := range requests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			request := httptest.NewRequest("GET", "/me", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code == http.StatusOK {
				successCount++
			} else if response.Code == http.StatusTooManyRequests {
				errorCount++
			}
		}(i)
	}

	wg.Wait()

	timeEnd := time.Now()

	testCase.Logf("Success Count: %d", successCount)
	testCase.Logf("Error Count: %d", errorCount)
	testCase.Logf("Time Elapsed: %v", timeEnd.Sub(timeStart))
	assert.Equal(testCase, numRequestsPerRoute, successCount)
}

func TestMaximumRequestInDifferentRoutesUsingSameMiddleware(testCase *testing.T) {
	router := SetupRouter()

	// Rate limiter configuration
	WindowCapacity := 5
	WindowLen := 1 * time.Second

	// Fake requests
	successCount := 0
	errorCount := 0
	numRequestsPerRoute := 5
	numRoutes := 2
	requests := make([]int, numRequestsPerRoute)

	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(&ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip_test",
		Option: ratelimiter.RateLimiterOption{
			Limit:  1,
			Burst:  WindowCapacity,
			Window: WindowLen,
		},
	})
	router.GET("/ping", rateLimiterMiddleware, func(ctx *gin.Context) {
		ratelimiter.StatusOK(ctx, gin.H{"message": "pong"})
	})

	router.GET("/me", rateLimiterMiddleware, func(ctx *gin.Context) {
		ratelimiter.StatusOK(ctx, gin.H{"message": "hello world"})
	})

	timeStart := time.Now()

	for i := range requests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			request := httptest.NewRequest("GET", "/me", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code == http.StatusOK {
				successCount++
			} else if response.Code == http.StatusTooManyRequests {
				errorCount++
			}
		}(i)
	}

	for i := range requests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			request := httptest.NewRequest("GET", "/ping", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code == http.StatusOK {
				successCount++
			} else if response.Code == http.StatusTooManyRequests {
				errorCount++
			}
		}(i)
	}

	wg.Wait()

	timeEnd := time.Now()

	testCase.Logf("Success Count: %d", successCount)
	testCase.Logf("Error Count: %d", errorCount)
	testCase.Logf("Time Elapsed: %v", timeEnd.Sub(timeStart))
	assert.Equal(testCase, numRequestsPerRoute*numRoutes, successCount+errorCount)
}
