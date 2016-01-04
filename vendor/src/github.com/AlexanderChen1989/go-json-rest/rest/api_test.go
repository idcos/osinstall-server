package rest

import (
	"testing"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"golang.org/x/net/context"
)

func TestAPINoAppNoMiddleware(t *testing.T) {

	api := NewAPI()
	if api == nil {
		t.Fatal("API object must be instantiated")
	}

	handler := api.MakeHandler()
	if handler == nil {
		t.Fatal("the http.Handler must be have been create")
	}

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/", nil))
	recorded.CodeIs(200)
}

func TestAPISimpleAppNoMiddleware(t *testing.T) {

	api := NewAPI()
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	handler := api.MakeHandler()
	if handler == nil {
		t.Fatal("the http.Handler must be have been create")
	}

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
	recorded.BodyIs(`{"Id":"123"}`)
}

func TestDevStack(t *testing.T) {

	api := NewAPI()
	api.Use(DefaultDevStack...)
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	handler := api.MakeHandler()
	if handler == nil {
		t.Fatal("the http.Handler must be have been create")
	}

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
	recorded.BodyIs("{\n  \"Id\": \"123\"\n}")
}

func TestProdStack(t *testing.T) {

	api := NewAPI()
	api.Use(DefaultProdStack...)
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	handler := api.MakeHandler()
	if handler == nil {
		t.Fatal("the http.Handler must be have been create")
	}

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
	recorded.ContentEncodingIsGzip()
}

func TestCommonStack(t *testing.T) {

	api := NewAPI()
	api.Use(DefaultCommonStack...)
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	handler := api.MakeHandler()
	if handler == nil {
		t.Fatal("the http.Handler must be have been create")
	}

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()
	recorded.BodyIs(`{"Id":"123"}`)
}
