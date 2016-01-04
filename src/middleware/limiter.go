package middleware

import (
	"logger"
	"time"

	"github.com/AlexanderChen1989/go-json-rest/rest"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"golang.org/x/net/context"
)

// LimiterMiddleware cancel context when timeout
type LimiterMiddleware struct {
	logger  logger.Logger
	limiter *config.Limiter
}

// Limiter create timeout middleware with duration
func Limiter(logger logger.Logger, limiter *config.Limiter) *LimiterMiddleware {
	return &LimiterMiddleware{logger, limiter}
}

// NewLimiterMiddleware create a simple limiter
func NewLimiterMiddleware(logger logger.Logger, max int64, ttl time.Duration) *LimiterMiddleware {
	return Limiter(logger, tollbooth.NewLimiter(max, ttl))
}

// MiddlewareFunc makes LimiterMiddleware implement the Middleware interface.
func (mw *LimiterMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
		httpError := tollbooth.LimitByRequest(mw.limiter, r.Request)
		if httpError != nil {
			mw.logger.Warnf("%s\n", httpError.Message)
			w.WriteHeader(httpError.StatusCode)
			w.WriteJSON(map[string]string{"status": "error", "msg": httpError.Message})
			return
		}

		// There's no rate-limit error, serve the next handler.
		h(ctx, w, r)
	}
}
