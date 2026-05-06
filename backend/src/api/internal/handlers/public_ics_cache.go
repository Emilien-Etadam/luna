package handlers

import (
	"sync"
	"time"
)

const publicICSCacheTTL = 5 * time.Minute

// In-memory cache for the generated public ICS body (single-instance deployments).
var publicICSBodyCache = struct {
	mu      sync.RWMutex
	body    []byte
	builtAt time.Time
	valid   bool
}{}

func getCachedPublicICSBody(now time.Time) ([]byte, bool) {
	publicICSBodyCache.mu.RLock()
	defer publicICSBodyCache.mu.RUnlock()
	if !publicICSBodyCache.valid {
		return nil, false
	}
	if now.Sub(publicICSBodyCache.builtAt) >= publicICSCacheTTL {
		return nil, false
	}
	return publicICSBodyCache.body, true
}

func setCachedPublicICSBody(body []byte, now time.Time) {
	publicICSBodyCache.mu.Lock()
	defer publicICSBodyCache.mu.Unlock()
	publicICSBodyCache.body = body
	publicICSBodyCache.builtAt = now
	publicICSBodyCache.valid = true
}

// InvalidatePublicICSCache clears the feed cache (admin toggle or future hooks).
func InvalidatePublicICSCache() {
	publicICSBodyCache.mu.Lock()
	defer publicICSBodyCache.mu.Unlock()
	publicICSBodyCache.valid = false
	publicICSBodyCache.body = nil
}
