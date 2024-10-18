package main

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

func Timeout(duration time.Duration) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			ctx, cancel := context.WithTimeout(ctx, duration)
			defer cancel()
			req = req.WithContext(ctx)
			h.ServeHTTP(rw, req)
		})
	}
}

func simulateOperation(ctx context.Context) (string, error) {
	wait := time.Duration(rand.Intn(200)) * time.Millisecond
	select {
	case <-time.After(wait):
		return "Done!", nil
	case <-ctx.Done():
		return "Too slow!", ctx.Err()
	}
}

func sleepyHandler(rw http.ResponseWriter, req *http.Request) {
	message, err := simulateOperation(req.Context())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			rw.WriteHeader(http.StatusGatewayTimeout)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		rw.WriteHeader(http.StatusOK)
	}
	rw.Write([]byte(message))
}

func main() {
	timeout := Timeout(100 * time.Millisecond)
	server := http.Server{
		Handler: timeout(http.HandlerFunc(sleepyHandler)),
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
