package api

import middleware "luna-backend/api/internal"

// CleanStalePublicMinuteLimiter garbage-collects entries in the public minute rate limiter (cron).
func CleanStalePublicMinuteLimiter() {
	middleware.CleanStalePublicMinuteLimiter()
}
