package requestid

import (
	"context"
	"jamdrop/logger"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type requestIDKey struct{}

var requestIDk requestIDKey

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(setToContext(r.Context(), generateRequestID())))
	})
}

func Log(r *http.Request) logger.Fieldser {
	return LogContext(r.Context())
}

func LogContext(ctx context.Context) logger.Fieldser {
	requestID := FromContext(ctx)
	if requestID == "" {
		return logger.Fields{}
	}

	return logger.Fields{"request_id": requestID}
}

func FromRequest(r *http.Request) string {
	return FromContext(r.Context())
}

func FromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDk).(string)
	if !ok {
		return ""
	}

	return requestID
}

func setToContext(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, requestIDk, requestID)
}

func generateRequestID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	const length = 16

	runes := make([]rune, length)
	for i := 0; i < length; i++ {
		runes[i] = rune(alphabet[rand.Intn(len(alphabet))])
	}

	return string(runes)
}
