package server

import (
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
			recorder := newLoggerRecorder(w)

			before := time.Now()
			next.ServeHTTP(recorder, r)
			elapsed := time.Since(before)

			requestSize := r.ContentLength
			if requestSize < 0 {
				requestSize = 0
			}

			logger.Printf("%s %s %dB | %d %s %dB | %v\n",
				// Request
				r.Method,
				r.URL.Path,
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
