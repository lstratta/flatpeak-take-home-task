package server

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s %s", r.Method, r.URL.Path, r.Host, time.Since(start))

		next.ServeHTTP(w, r)

	})
}
