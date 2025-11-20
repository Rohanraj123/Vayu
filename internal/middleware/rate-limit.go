package middleware

import (
	"net/http"
	"time"

	"github.com/Rohanraj123/vayu/internal/config"
)

// RateLimiter object is used for dynamic runtime tracking
type rateLimiter struct {
	tokens     int
	capacity   int     // burst
	refillRate float64 // requests_per_minute
	lastRefill time.Duration
}

// new object generator
func NewRateLimiter(
	tokens int,
	capacity int,
	refillRate float64,
	lastRefill time.Duration,
) *rateLimiter {
	return &rateLimiter{
		tokens:     tokens,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: lastRefill,
	}
}

// rate-limit handler func
func RateLimitMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

// admin will configure the rate_limitting using config.yaml file
// Then when ever it generates the API_KEY, it stores the in-memory map for API_KEY and rateLimitConfig
// Once the user request with API_KEY then it extracts the API_KEY and then check whether tokenBucket is
// available or not. If available, it then decreases the value by 1 and allows the request and if not
// then decline the request.

// it also keep on checking the time and refils the bucket again and again upto full capacity.
