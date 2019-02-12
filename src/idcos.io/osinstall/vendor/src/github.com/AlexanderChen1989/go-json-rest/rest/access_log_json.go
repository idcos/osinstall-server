package rest

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
)

// AccessLogJSONMiddleware produces the access log with records written as JSON. This middleware
// depends on TimerMiddleware and RecorderMiddleware that must be in the wrapped middlewares. It
// also uses request.Env["REMOTE_USER"].(string) set by the auth middlewares.
type AccessLogJSONMiddleware struct {

	// Logger points to the logger object used by this middleware, it defaults to
	// log.New(os.Stderr, "", 0).
	Logger *log.Logger
}

// MiddlewareFunc makes AccessLogJSONMiddleware implement the Middleware interface.
func (mw *AccessLogJSONMiddleware) MiddlewareFunc(h HandlerFunc) HandlerFunc {

	// set the default Logger
	if mw.Logger == nil {
		mw.Logger = log.New(os.Stderr, "", 0)
	}

	return func(ctx context.Context, w ResponseWriter, r *Request) {

		// call the handler
		h(ctx, w, r)

		mw.Logger.Printf("%s", makeAccessLogJSONRecord(ctx, r).asJSON())
	}
}

// AccessLogJSONRecord is the data structure used by AccessLogJSONMiddleware to create the JSON
// records. (Public for documentation only, no public method uses it)
type AccessLogJSONRecord struct {
	Timestamp    *time.Time
	StatusCode   int
	ResponseTime *time.Duration
	HTTPMethod   string
	RequestURI   string
	RemoteUser   string
	UserAgent    string
}

func makeAccessLogJSONRecord(ctx context.Context, r *Request) *AccessLogJSONRecord {
	env := EnvFromContext(ctx)
	var timestamp *time.Time
	if env["START_TIME"] != nil {
		timestamp = env["START_TIME"].(*time.Time)
	}

	var statusCode int
	if env["STATUS_CODE"] != nil {
		statusCode = env["STATUS_CODE"].(int)
	}

	var responseTime *time.Duration
	if env["ELAPSED_TIME"] != nil {
		responseTime = env["ELAPSED_TIME"].(*time.Duration)
	}

	var remoteUser string
	if env["REMOTE_USER"] != nil {
		remoteUser = env["REMOTE_USER"].(string)
	}

	return &AccessLogJSONRecord{
		Timestamp:    timestamp,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		HTTPMethod:   r.Method,
		RequestURI:   r.URL.RequestURI(),
		RemoteUser:   remoteUser,
		UserAgent:    r.UserAgent(),
	}
}

func (r *AccessLogJSONRecord) asJSON() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return b
}
