package middleware

import (
	"testing"
	"time"
)

func TestRateLimiterStopIsIdempotent(t *testing.T) {
	rl := NewRateLimiter(10, time.Minute)
	rl.Stop()
	rl.Stop() // 不得因重复 close 而 panic
}
