package rest

import (
	"testing"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"golang.org/x/net/context"
)

func TestContentTypeCheckerMiddleware(t *testing.T) {

	api := NewAPI()

	// the middleware to test
	api.Use(&ContentTypeCheckerMiddleware{})

	// a simple app
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	// wrap all
	handler := api.MakeHandler()

	// no payload, no content length, no check
	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/", nil))
	recorded.CodeIs(200)

	// JSON payload with correct content type
	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/", map[string]string{"Id": "123"}))
	recorded.CodeIs(200)

	// JSON payload with correct content type specifying the utf-8 charset
	req := test.MakeSimpleRequest("POST", "http://localhost/", map[string]string{"Id": "123"})
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	recorded = test.RunRequest(t, handler, req)
	recorded.CodeIs(200)

	// JSON payload with incorrect content type
	req = test.MakeSimpleRequest("POST", "http://localhost/", map[string]string{"Id": "123"})
	req.Header.Set("Content-Type", "text/x-json")
	recorded = test.RunRequest(t, handler, req)
	recorded.CodeIs(415)
}
