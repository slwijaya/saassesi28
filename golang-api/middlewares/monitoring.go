package middlewares

import (
	"log"
	"net/http"
	"time"
)

// CustomResponseWriter untuk mencatat status HTTP
type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *CustomResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// MonitoringMiddleware mencatat waktu respons API dan status kode
func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		crw := &CustomResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(crw, r)

		duration := time.Since(startTime)
		log.Printf("[MONITORING] %s %s - %d - %v", r.Method, r.URL.Path, crw.StatusCode, duration)
	})
}
