package middleware

import (
	"config"

	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
)

// ctxConfigKey 注入的*config.Config对应的查询Key
var ctxConfigKey uint8

// ConfigFromContext 从ctx中获取model.Repo
func ConfigFromContext(ctx context.Context) (*config.Config, bool) {
	conf, ok := ctx.Value(&ctxConfigKey).(*config.Config)
	return conf, ok
}

// InjectConfig 注入*config.Config
func InjectConfig(conf *config.Config) rest.Middleware {
	fn := func(h rest.HandlerFunc) rest.HandlerFunc {
		return func(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
			ctx = context.WithValue(ctx, &ctxConfigKey, conf)
			h(ctx, w, r)
		}
	}
	return rest.MiddlewareSimple(fn)
}
