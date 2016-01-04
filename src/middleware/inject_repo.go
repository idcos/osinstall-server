package middleware

import (
	"model"

	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
)

// ctxRepoKey 注入的model.Repo对应的查询Key
var ctxRepoKey uint8

// RepoFromContext 从ctx中获取model.Repo
func RepoFromContext(ctx context.Context) (model.Repo, bool) {
	repo, ok := ctx.Value(&ctxRepoKey).(model.Repo)
	return repo, ok
}

// InjectRepo 注入model.Repo
func InjectRepo(repo model.Repo) rest.Middleware {
	fn := func(h rest.HandlerFunc) rest.HandlerFunc {
		return func(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
			ctx = context.WithValue(ctx, &ctxRepoKey, repo)
			h(ctx, w, r)
		}
	}
	return rest.MiddlewareSimple(fn)
}
