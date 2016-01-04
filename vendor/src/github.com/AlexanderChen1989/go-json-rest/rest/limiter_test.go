package rest

import (
	"testing"
	"time"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"github.com/didip/tollbooth"
	"golang.org/x/net/context"
)

func TestLimiterMiddleware(t *testing.T) {
	api := NewAPI()

	// the middleware to test
	api.Use(Limiter(tollbooth.NewLimiter(1, time.Second)))

	// a simple app
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	// wrap all
	handler := api.MakeHandler()

	for i := 0; i < 100; i++ {
		req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
		req.RemoteAddr = "127.0.0.1:45344"
		recorded := test.RunRequest(t, handler, req)
		if i > 1 {
			recorded.CodeIs(429)
			recorded.ContentTypeIsJSON()
		}
	}
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
		req.RemoteAddr = "127.0.0.1:45344"
		recorded := test.RunRequest(t, handler, req)
		recorded.CodeIs(200)
		recorded.ContentTypeIsJSON()
	}
}
