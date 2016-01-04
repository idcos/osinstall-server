package rest

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"golang.org/x/net/context"
)

func TestAccessLogJSONMiddleware(t *testing.T) {

	api := NewAPI()

	// the middlewares stack
	buffer := bytes.NewBufferString("")
	api.Use(&AccessLogJSONMiddleware{
		Logger: log.New(buffer, "", 0),
	})
	api.Use(&TimerMiddleware{})
	api.Use(&RecorderMiddleware{})

	// a simple app
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJSON(map[string]string{"Id": "123"})
	}))

	// wrap all
	handler := api.MakeHandler()

	req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	recorded := test.RunRequest(t, handler, req)
	recorded.CodeIs(200)
	recorded.ContentTypeIsJSON()

	// log tests
	decoded := &AccessLogJSONRecord{}
	err := json.Unmarshal(buffer.Bytes(), decoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.StatusCode != 200 {
		t.Errorf("StatusCode 200 expected, got %d", decoded.StatusCode)
	}
	if decoded.RequestURI != "/" {
		t.Errorf("RequestURI / expected, got %s", decoded.RequestURI)
	}
	if decoded.HTTPMethod != "GET" {
		t.Errorf("HTTPMethod GET expected, got %s", decoded.HTTPMethod)
	}
}
