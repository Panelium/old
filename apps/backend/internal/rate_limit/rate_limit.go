package rate_limit

import (
	"sync"
	"time"
)

type visitor struct {
	lastSeen time.Time
	requests int
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	v, exists := rl.visitors[ip]
	now := time.Now()
	if !exists || now.Sub(v.lastSeen) > rl.window {
		rl.visitors[ip] = &visitor{lastSeen: now, requests: 1}
		return true
	}
	if v.requests >= rl.limit {
		return false
	}
	v.requests++
	v.lastSeen = now
	return true
}
