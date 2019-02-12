package rest

import (
	"testing"
	"time"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"golang.org/x/net/context"
)

func TestTimerMiddleware(t *testing.T) {

	api := NewAPI()

	// a middleware carrying the Env tests
	api.Use(MiddlewareSimple(func(handler HandlerFunc) HandlerFunc {
		return func(ctx context.Context, w ResponseWriter, r *Request) {

			handler(ctx, w, r)

			env := EnvFromContext(ctx)

			if env["ELAPSED_TIME"] == nil {
				t.Error("ELAPSED_TIME is nil")
			}
			elapsedTime := env["ELAPSED_TIME"].(*time.Duration)
			if elapsedTime.Nanoseconds() <= 0 {
				t.Errorf(
					"ELAPSED_TIME is expected to be at least 1 nanosecond %d",
					elapsedTime.Nanoseconds(),
				)
			}

			if env["START_TIME"] == nil {
				t.Error("START_TIME is nil")
			}
			start := env["START_TIME"].(*time.Time)
			if start.After(time.Now()) {
				t.Errorf(
					"START_TIME is expected to be in the past %s",
					start.String(),
				)
			}
		}
	}))

	// the middleware to test
	api.Use(&TimerMiddleware{})

	// a simple app
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	// wrap all
	handler := api.MakeHandler()

	req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	recorded := test.RunRequest(t, handler, req)
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
}
