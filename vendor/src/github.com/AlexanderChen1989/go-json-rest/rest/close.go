package rest

import (
	"net/http"

	"golang.org/x/net/context"
)

// CloseMiddleware cancel context when timeout
type CloseMiddleware struct{}

// MiddlewareFunc makes CloseMiddleware implement the Middleware interface.
func (mw *CloseMiddleware) MiddlewareFunc(h HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w ResponseWriter, r *Request) {
		// Cancel the context if the client closes the connection
		if wcn, ok := w.(http.CloseNotifier); ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			defer cancel()

			notify := wcn.CloseNotify()
			go func() {
				<-notify
				cancel()
			}()
		}

		h(ctx, w, r)
	}
}
