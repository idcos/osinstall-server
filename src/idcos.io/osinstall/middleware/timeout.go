package middleware

import (
	"time"

	"github.com/AlexanderChen1989/go-json-rest/rest"

	"golang.org/x/net/context"
)

// TimeoutMiddleware cancel context when timeout
type TimeoutMiddleware struct {
	timeout time.Duration
}

// NewTimeoutMiddleware create timeout middleware with duration
func NewTimeoutMiddleware(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{timeout}
}

// MiddlewareFunc makes TimeoutMiddleware implement the Middleware interface.
func (mw *TimeoutMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
		ctx, _ = context.WithTimeout(ctx, mw.timeout)
		h(ctx, w, r)
	}
}
