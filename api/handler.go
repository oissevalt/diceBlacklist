package api

import (
	"golang.org/x/time/rate"
	"net/http"
)

var limiter = rate.NewLimiter(3, 12)

func Limit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			f(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	}
}
