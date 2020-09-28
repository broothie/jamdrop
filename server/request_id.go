package server

import (
	"context"
	"math/rand"
	"net/http"
)

type requestIDKey struct{}

var requestIDK requestIDKey

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDK, newID())))
	})
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDK).(string)
	return requestID, ok
}

func newID() string {
	const (
		length   = 8
		alphabet = "abcdefghijklmnopqrstuvwxyz"
	)

	runes := make([]rune, length)
	for i := 0; i < length; i++ {
		runes[i] = rune(alphabet[rand.Intn(len(alphabet))])
	}

	return string(runes)
}
