package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/Rohanraj123/vayu/internal/config"
)

// RateLimiter object is used for dynamic runtime tracking
type RateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64 // burst
	refillRate float64 // tokens_per_second
	lastRefill time.Time
}

// create a new rate-limiter
func NewRateLimiter(
	cfg *config.RateLimitConfig,
) *RateLimiter {
	return &RateLimiter{
		tokens:     float64(cfg.Burst),
		capacity:   float64(cfg.Burst),
		refillRate: float64(cfg.RequestPerMinute) / 60.0,
		lastRefill: time.Now(),
	}
}

// Store all API-KEY limiters
var limitStore = struct {
	sync.RWMutex
	store map[string]*RateLimiter
}{store: make(map[string]*RateLimiter)}

// rate-limit handler func
func RateLimitMiddleware(cfg *config.RateLimitConfig, next http.Handler) http.Handler {
	if !cfg.Enabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(w, "missing API KEY for rate-limiting", http.StatusUnauthorized)
			return
		}

		limiter := getLimiter(apiKey, *cfg)
		if !limiter.Allow() {
			http.Error(w, "TOO MANY REQUESTS - rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Get or Create RateLimiter API-KEY
func getLimiter(apiKey string, cfg config.RateLimitConfig) *RateLimiter {
	limitStore.RLock()
	limiter := limitStore.store[apiKey]
	limitStore.RUnlock()

	if limiter != nil {
		return limiter
	}

	// Create a new limiter
	limitStore.Lock()
	defer limitStore.Unlock()

	limiter = NewRateLimiter(&cfg)
	limitStore.store[apiKey] = limiter
	return limiter
}

// consume token
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	rl.lastRefill = now

	// Check if a token is available
	if rl.tokens >= 1 {
		rl.tokens -= 1
		return true
	}

	return false
}
