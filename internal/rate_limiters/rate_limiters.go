package rate_limiters

import (
	"sync"

	"github.com/sergey-kurenkov/test_number_requests/internal/rate_limiter"
)

type RateLimiters struct {
	mtx sync.RWMutex
	rls map[string]*rate_limiter.RateLimiter
}

func NewRateLimiters() *RateLimiters {
	return &RateLimiters{
		rls: make(map[string]*rate_limiter.RateLimiter),
	}
}

func (c *RateLimiters) GetRateLimiter(ipAddress string) *rate_limiter.RateLimiter {
	rl, ok := func() (*rate_limiter.RateLimiter, bool) {
		c.mtx.RLock()
		defer c.mtx.RUnlock()
		rl, ok := c.rls[ipAddress]
		return rl, ok
	}()

	if ok {
		return rl
	}

	rl = func() *rate_limiter.RateLimiter {
		c.mtx.Lock()
		defer c.mtx.Unlock()
		rl, ok := c.rls[ipAddress]
		if ok {
			return rl
		}

		rl = rate_limiter.NewRateLimiter(5)
		c.rls[ipAddress] = rl

		return rl
	}()

	return rl
}
