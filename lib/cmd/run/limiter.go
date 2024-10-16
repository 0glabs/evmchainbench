package run

import (
	"fmt"
	"sync"
)

type RateLimiter struct {
	mutex     sync.Mutex
	remaining int
}

func NewRateLimiter(maxRequests int) *RateLimiter {
	return &RateLimiter{
		remaining: maxRequests,
	}
}

func (rl *RateLimiter) AllowRequest() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.remaining > 0 {
		rl.remaining--
		return true
	}
	return false
}

func (rl *RateLimiter) IncreaseLimit(newCapacity int) {
	rl.mutex.Lock()
	rl.remaining += newCapacity
	rl.mutex.Unlock()
	fmt.Println("Remaining requests:", rl.remaining)
}
