package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	rate      float64   // tokens per second
	burst     int       // max tokens
	tokens    float64   // current tokens
	lastCheck time.Time // last refill timestamp
	mu        sync.Mutex
}

type Limiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
}

func NewInMemoryLimiter() *Limiter {
	return &Limiter{
		buckets: make(map[string]*TokenBucket),
	}
}

func (l *Limiter) SetRate(key string, rate float64, burst int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buckets[key] = &TokenBucket{
		rate:      rate,
		burst:     burst,
		tokens:    float64(burst),
		lastCheck: time.Now(),
	}

}

func (l *Limiter) Allow(key string) bool {
	l.mu.RLock()
	bucket, exists := l.buckets[key]
	l.mu.RUnlock()
	if !exists {
		return true // default allow if no rate set
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastCheck).Seconds()
	bucket.tokens += elapsed * bucket.rate
	if bucket.tokens > float64(bucket.burst) {
		bucket.tokens = float64(bucket.burst)
	}
	bucket.lastCheck = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}
	return false
}
