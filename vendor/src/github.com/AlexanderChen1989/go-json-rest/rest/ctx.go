package rest

import "golang.org/x/net/context"

const ctxPathParams = "PATH_PARAMS"
const ctxEnv = "ENV"

func contextWithEnv() context.Context {
	return context.WithValue(
		context.Background(),
		ctxEnv,
		&map[string]interface{}{},
	)
}

// PathParamFromContext fetch PathParam from context
func PathParamFromContext(ctx context.Context) map[string]string {
	return *(ctx.Value(ctxPathParams).(*map[string]string))
}

// EnvFromContext fetch Env from context
func EnvFromContext(ctx context.Context) map[string]interface{} {
	return *(ctx.Value(ctxEnv).(*map[string]interface{}))
}
