package rest

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"golang.org/x/net/context"
)

// LimiterMiddleware cancel context when timeout
type LimiterMiddleware struct {
	limiter *config.Limiter
}

// Limiter create timeout middleware with duration
func Limiter(limiter *config.Limiter) *LimiterMiddleware {
	return &LimiterMiddleware{limiter}
}

// SimpleLimiter create a simple limiter
func SimpleLimiter(max int64, ttl time.Duration) *LimiterMiddleware {
	return Limiter(tollbooth.NewLimiter(max, ttl))
}

// MiddlewareFunc makes LimiterMiddleware implement the Middleware interface.
func (mw *LimiterMiddleware) MiddlewareFunc(h HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w ResponseWriter, r *Request) {
		httpError := tollbooth.LimitByRequest(mw.limiter, r.Request)
		if httpError != nil {
			w.WriteHeader(httpError.StatusCode)
			w.WriteJSON(map[string]string{"status": "error", "msg": httpError.Message})
			return
		}

		// There's no rate-limit error, serve the next handler.
		h(ctx, w, r)
	}
}
