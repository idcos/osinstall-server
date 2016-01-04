package middleware

import (
	"logger"
	"net/http"

	"github.com/AlexanderChen1989/go-json-rest/rest"

	"golang.org/x/net/context"
)

// CloseMiddleware cancel context when timeout
type CloseMiddleware struct {
	logger logger.Logger
}

// NewCloseMiddleware create middleware
func NewCloseMiddleware(log logger.Logger) *CloseMiddleware {
	return &CloseMiddleware{log}
}

// MiddlewareFunc makes CloseMiddleware implement the Middleware interface.
func (mw *CloseMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
		// Cancel the context if the client closes the connection
		if wcn, ok := w.(http.CloseNotifier); ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			defer cancel()

			notify := wcn.CloseNotify()
			go func() {
				<-notify
				mw.logger.Warn("Remote closed, cancel context.\n")
				cancel()
			}()
		}

		h(ctx, w, r)
	}
}
