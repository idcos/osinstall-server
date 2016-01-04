package middleware

import (
	"logger"

	"github.com/AlexanderChen1989/go-json-rest/rest"

	"golang.org/x/net/context"
)

// ctxLoggerKey 注入的logger.Logger对应的查询Key
var ctxLoggerKey uint8

// LoggerFromContext 从ctx中获取model.Repo
func LoggerFromContext(ctx context.Context) (logger.Logger, bool) {
	log, ok := ctx.Value(&ctxLoggerKey).(logger.Logger)
	return log, ok
}

// InjectLogger 注入logger.Logger
func InjectLogger(logger logger.Logger) rest.Middleware {
	fn := func(h rest.HandlerFunc) rest.HandlerFunc {
		return func(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
			ctx = context.WithValue(ctx, &ctxLoggerKey, logger)
			h(ctx, w, r)
		}
	}
	return rest.MiddlewareSimple(fn)
}
