package rest

import (
	"net/http"

	"golang.org/x/net/context"
)

// API defines a stack of Middlewares and an App.
type API struct {
	stack []Middleware
	app   App
}

// NewAPI makes a new API object. The Middleware stack is empty, and the App is nil.
func NewAPI() *API {
	return &API{
		stack: []Middleware{},
		app:   nil,
	}
}

// Use pushes one or multiple middlewares to the stack for middlewares
// maintained in the API object.
func (api *API) Use(middlewares ...Middleware) {
	api.stack = append(api.stack, middlewares...)
}

// SetApp sets the App in the API object.
func (api *API) SetApp(app App) {
	api.app = app
}

// MakeHandler wraps all the Middlewares of the stack and the App together, and returns an
// http.Handler ready to be used. If the Middleware stack is empty the App is used directly. If the
// App is nil, a HandlerFunc that does nothing is used instead.
func (api *API) MakeHandler() http.Handler {
	var appFunc HandlerFunc
	if api.app != nil {
		appFunc = api.app.AppFunc()
	} else {
		appFunc = func(ctx context.Context, w ResponseWriter, r *Request) {}
	}
	return http.HandlerFunc(
		adapterFunc(
			WrapMiddlewares(api.stack, appFunc),
		),
	)
}

// Defines a stack of middlewares convenient for development. Among other things:
// console friendly logging, JSON indentation, error stack strace in the response.
var DefaultDevStack = []Middleware{
	&AccessLogApacheMiddleware{},
	&TimerMiddleware{},
	&RecorderMiddleware{},
	&PoweredByMiddleware{},
	&RecoverMiddleware{
		EnableResponseStackTrace: true,
	},
	&JSONIndentMiddleware{},
	&ContentTypeCheckerMiddleware{},
}

// Defines a stack of middlewares convenient for production. Among other things:
// Apache CombinedLogFormat logging, gzip compression.
var DefaultProdStack = []Middleware{
	&AccessLogApacheMiddleware{
		Format: CombinedLogFormat,
	},
	&TimerMiddleware{},
	&RecorderMiddleware{},
	&PoweredByMiddleware{},
	&RecoverMiddleware{},
	&GzipMiddleware{},
	&ContentTypeCheckerMiddleware{},
}

// Defines a stack of middlewares that should be common to most of the middleware stacks.
var DefaultCommonStack = []Middleware{
	&TimerMiddleware{},
	&RecorderMiddleware{},
	&PoweredByMiddleware{},
	&RecoverMiddleware{},
}
