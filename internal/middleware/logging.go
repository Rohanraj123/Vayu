package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware wraps an HTTP handler and logs request + response info
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: 200} // default 200

		// Log BEFORE processing
		log.Printf("[START] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler (e.g. reverseProxy)
		next.ServeHTTP(rw, r)

		// Log AFTER processing
		duration := time.Since(start)
		log.Printf("[END] %s %s %d completed in %v", r.Method, r.URL.Path, rw.statusCode, duration)
	})
}
