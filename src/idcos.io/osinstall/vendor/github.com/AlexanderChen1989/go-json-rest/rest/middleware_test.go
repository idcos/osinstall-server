package rest

import (
	"testing"

	"golang.org/x/net/context"
)

type testMiddleware struct {
	name string
}

func (mw *testMiddleware) MiddlewareFunc(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w ResponseWriter, r *Request) {
		env := EnvFromContext(ctx)
		if env["BEFORE"] == nil {
			env["BEFORE"] = mw.name
		} else {
			env["BEFORE"] = env["BEFORE"].(string) + mw.name
		}
		handler(ctx, w, r)
		if env["AFTER"] == nil {
			env["AFTER"] = mw.name
		} else {
			env["AFTER"] = env["AFTER"].(string) + mw.name
		}
	}
}

func TestWrapMiddlewares(t *testing.T) {

	a := &testMiddleware{"A"}
	b := &testMiddleware{"B"}
	c := &testMiddleware{"C"}

	app := func(ctx context.Context, w ResponseWriter, r *Request) {
		// do nothing
	}

	handlerFunc := WrapMiddlewares([]Middleware{a, b, c}, app)

	ctx := contextWithEnv()
	// fake request
	r := &Request{}

	handlerFunc(ctx, nil, r)

	env := EnvFromContext(ctx)
	before := env["BEFORE"].(string)
	if before != "ABC" {
		t.Error("middleware executed in the wrong order, expected ABC")
	}

	after := env["AFTER"].(string)
	if after != "CBA" {
		t.Error("middleware executed in the wrong order, expected CBA")
	}
}
