package middlewares

import (
	"log"
	"net/http"
	"time"
)

type LogMiddleware struct {
	next http.Handler
}

func Log(h http.Handler) *LogMiddleware {
	return &LogMiddleware{next: h}
}

func (m *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startAt := time.Now()
	m.next.ServeHTTP(w, r)

	log.Printf("[http] %s %s - %d ms",
		r.Method,
		r.URL.Path,
		time.Since(startAt).Truncate(time.Millisecond)/time.Millisecond)
}
