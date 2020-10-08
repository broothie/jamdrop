package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type loggerRecorder struct {
	http.ResponseWriter
	status     int
	bodyLength int
}

func newLoggerRecorder(w http.ResponseWriter) *loggerRecorder {
	return &loggerRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (lr *loggerRecorder) WriteHeader(status int) {
	lr.status = status
	lr.ResponseWriter.WriteHeader(status)
}

func (lr *loggerRecorder) Write(body []byte) (int, error) {
	lr.bodyLength = len(body)
	return lr.ResponseWriter.Write(body)
}

func ApplyLoggerMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return LoggerMiddleware(logger)(next)
}

func LoggerMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestSize := r.ContentLength
			if requestSize < 0 {
				requestSize = 0
			}

			query := ""
			if len(r.URL.RawQuery) > 0 {
				query = fmt.Sprintf("?%s", r.URL.RawQuery)
			}

			// Make request
			recorder := newLoggerRecorder(w)
			before := time.Now()
			next.ServeHTTP(recorder, r)
			elapsed := time.Since(before)

			// Log after
			logger.Printf("%s %s%s %dB | %d %s %dB | %v\n",
				// Request
				r.Method,
				r.URL.Path,
				query,
				requestSize,
				//Response
				recorder.status,
				http.StatusText(recorder.status),
				recorder.bodyLength,
				// Timing
				elapsed,
			)
		})
	}
}