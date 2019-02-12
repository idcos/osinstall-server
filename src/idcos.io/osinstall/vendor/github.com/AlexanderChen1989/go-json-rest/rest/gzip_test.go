package rest

import (
	"testing"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"golang.org/x/net/context"
)

func TestGzipEnabled(t *testing.T) {

	api := NewAPI()

	// the middleware to test
	api.Use(&GzipMiddleware{})

	// router app with success and error paths
	router, err := MakeRouter(
		Get("/ok", func(ctx context.Context, w ResponseWriter, r *Request) {
			w.WriteJSON(map[string]string{"Id": "123"})
		}),
		Get("/error", func(ctx context.Context, w ResponseWriter, r *Request) {
			Error(w, "gzipped error", 500)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	api.SetApp(router)

	// wrap all
	handler := api.MakeHandler()

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/ok", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
	recorded.ContentEncodingIsGzip()
	recorded.HeaderIs("Vary", "Accept-Encoding")

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/error", nil))
	recorded.CodeIs(500)
	recorded.ContentTypeIsJSON()
	recorded.ContentEncodingIsGzip()
	recorded.HeaderIs("Vary", "Accept-Encoding")
}

func TestGzipDisabled(t *testing.T) {

	api := NewAPI()

	// router app with success and error paths
	router, err := MakeRouter(
		Get("/ok", func(ctx context.Context, w ResponseWriter, r *Request) {
			w.WriteJSON(map[string]string{"Id": "123"})
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	api.SetApp(router)
	handler := api.MakeHandler()

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/ok", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
	recorded.HeaderIs("Content-Encoding", "")
	recorded.HeaderIs("Vary", "")
}
