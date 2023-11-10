package api

import (
	"diceBlacklist/db"
	"golang.org/x/time/rate"
	"net/http"
)

var limiter = rate.NewLimiter(3, 12)

func Limit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() || !db.Running.Load() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		} else {
			f(w, r)
		}
	}
}
