package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu sync.Mutex
	visitors map[string]int
	resetTime time.Duration
	limit int
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter{
	rl := &rateLimiter{
		visitors: make(map[string]int),
		resetTime: resetTime,
		limit: limit,
	}
	go rl.resetVisitorCount()
	return rl
}

func (rl *rateLimiter) resetVisitorCount(){
	time.Sleep(rl.resetTime)
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.visitors = make(map[string]int)
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler{

	fmt.Println("RateLimiter middlware")
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RateLimiter middlware being returned...")

		rl.mu.Lock()
		defer rl.mu.Unlock()

		visitorIp := r.RemoteAddr
		rl.visitors[visitorIp]++
		fmt.Printf("Visitor count from %v is %v\n", visitorIp, rl.visitors[visitorIp])

		if rl.visitors[visitorIp]>rl.limit {
			http.Error(w, "Too many request", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
		fmt.Println("RateLimiter middlware ends")
	})

}