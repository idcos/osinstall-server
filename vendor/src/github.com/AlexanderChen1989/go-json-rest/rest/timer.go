package rest

import (
	"time"

	"golang.org/x/net/context"
)

// TimerMiddleware computes the elapsed time spent during the execution of the wrapped handler.
// The result is available to the wrapping handlers as request.Env["ELAPSED_TIME"].(*time.Duration),
// and as request.Env["START_TIME"].(*time.Time)
type TimerMiddleware struct{}

// MiddlewareFunc makes TimerMiddleware implement the Middleware interface.
func (mw *TimerMiddleware) MiddlewareFunc(h HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w ResponseWriter, r *Request) {

		start := time.Now()
		env := EnvFromContext(ctx)
		env["START_TIME"] = &start

		// call the handler
		h(ctx, w, r)

		end := time.Now()
		elapsed := end.Sub(start)
		env["ELAPSED_TIME"] = &elapsed
	}
}
