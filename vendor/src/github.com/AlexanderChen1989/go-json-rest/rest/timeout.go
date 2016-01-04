package rest

import (
	"time"

	"golang.org/x/net/context"
)

// TimeoutMiddleware cancel context when timeout
type TimeoutMiddleware struct {
	timeout time.Duration
}

// Timeout create timeout middleware with duration
func Timeout(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{timeout}
}

// MiddlewareFunc makes TimeoutMiddleware implement the Middleware interface.
func (mw *TimeoutMiddleware) MiddlewareFunc(h HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w ResponseWriter, r *Request) {
		ctx, _ = context.WithTimeout(ctx, mw.timeout)
		h(ctx, w, r)
	}
}
