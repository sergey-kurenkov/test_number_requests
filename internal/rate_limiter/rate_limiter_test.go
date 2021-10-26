package rate_limiter

import (
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	rl := NewRateLimiter(2)

	l := rl.Filled()
	if l != 2 {
		t.Error(l)
	}

	rl.AddRequest()
	rl.AddRequest()

	l = rl.Filled()
	if l != 0 {
		t.Error(l)
	}

	rl.OnFinishRequest()
	l = rl.Filled()
	if l != 1 {
		t.Error(l)
	}

	rl.AddRequest()

	l = rl.Filled()
	if l != 0 {
		t.Error(l)
	}

	w := rl.Waiting()
	if w != 0 {
		t.Error(w)
	}

	go func() {
		rl.AddRequest()
	}()

	time.Sleep(1 * time.Second)

	w = rl.Waiting()
	if w != 1 {
		t.Error(w)
	}

	rl.OnFinishRequest()
	time.Sleep(1 * time.Millisecond)

	w = rl.Waiting()
	if w != 0 {
		t.Error(w)
	}

}
