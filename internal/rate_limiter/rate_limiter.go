package rate_limiter

import "sync/atomic"

type RateLimiter struct {
	ch chan struct{}
	waiting int64
}

func NewRateLimiter(capacity int) *RateLimiter {
	return &RateLimiter{
		ch: make(chan struct{}, capacity),
	}
}

func (rl *RateLimiter) Filled() int {
	c := cap(rl.ch)
	l := len(rl.ch)
	f := c - l
	return f
}

func (rl *RateLimiter) Waiting() int64 {
	res := atomic.LoadInt64(&rl.waiting)
	return res
}

func (rl *RateLimiter) AddRequest() {
	atomic.AddInt64(&rl.waiting, 1)
	rl.ch <- struct{}{}

	atomic.AddInt64(&rl.waiting, -1)
}

func (rl *RateLimiter) OnFinishRequest() {
	<- rl.ch
}