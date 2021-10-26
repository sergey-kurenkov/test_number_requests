package rate_limiters

import (
	"testing"
)

func TestRateLimiters(t *testing.T) {
	testIP1 := "127.0.0.1"
	testIP2 := "127.0.0.1"
	testIP3 := "192.168.1.1"

	rls := NewRateLimiters()

	rl1 := rls.GetRateLimiter(testIP1)
	rl2 := rls.GetRateLimiter(testIP2)
	if rl1 != rl2 {
		t.Fatal(rl1, rl2)
	}

	rl3 := rls.GetRateLimiter(testIP3)
	if rl1 == rl3 {
		t.Fatal(rl1, rl3)
	}
}
