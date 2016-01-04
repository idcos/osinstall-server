package rest

import (
	"testing"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"golang.org/x/net/context"
)

func TestRecorderMiddleware(t *testing.T) {

	api := NewAPI()

	// a middleware carrying the Env tests
	api.Use(MiddlewareSimple(func(handler HandlerFunc) HandlerFunc {
		return func(ctx context.Context, w ResponseWriter, r *Request) {

			handler(ctx, w, r)
			env := EnvFromContext(ctx)
			if env["STATUS_CODE"] == nil {
				t.Error("STATUS_CODE is nil")
			}
			statusCode := env["STATUS_CODE"].(int)
			if statusCode != 200 {
				t.Errorf("STATUS_CODE = 200 expected, got %d", statusCode)
			}

			if env["BYTES_WRITTEN"] == nil {
				t.Error("BYTES_WRITTEN is nil")
			}
			bytesWritten := env["BYTES_WRITTEN"].(int64)
			// '{"Id":"123"}' => 12 chars
			if bytesWritten != 12 {
				t.Errorf("BYTES_WRITTEN 12 expected, got %d", bytesWritten)
			}
		}
	}))

	// the middleware to test
	api.Use(&RecorderMiddleware{})

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

// See how many bytes are written when gzipping
func TestRecorderAndGzipMiddleware(t *testing.T) {

	api := NewAPI()

	// a middleware carrying the Env tests
	api.Use(MiddlewareSimple(func(handler HandlerFunc) HandlerFunc {
		return func(ctx context.Context, w ResponseWriter, r *Request) {

			handler(ctx, w, r)
			env := EnvFromContext(ctx)
			if env["BYTES_WRITTEN"] == nil {
				t.Error("BYTES_WRITTEN is nil")
			}
			bytesWritten := env["BYTES_WRITTEN"].(int64)
			// Yes, the gzipped version actually takes more space.
			if bytesWritten != 28 {
				t.Errorf("BYTES_WRITTEN 28 expected, got %d", bytesWritten)
			}
		}
	}))

	// the middlewares to test
	api.Use(&RecorderMiddleware{})
	api.Use(&GzipMiddleware{})

	// a simple app
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	// wrap all
	handler := api.MakeHandler()

	req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	// "Accept-Encoding", "gzip" is set by test.MakeSimpleRequest
	recorded := test.RunRequest(t, handler, req)
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
}
