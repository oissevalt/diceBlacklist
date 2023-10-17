package api

import (
	"golang.org/x/time/rate"
	"net/http"
)

var limiter *rate.Limiter

func init() {
	limiter = rate.NewLimiter(3, 12)
}

func Limit(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			f(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	}
}
