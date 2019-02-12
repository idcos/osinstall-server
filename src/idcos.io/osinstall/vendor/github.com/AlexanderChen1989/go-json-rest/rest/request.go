package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
	ErrJSONPayloadEmpty = errors.New("JSON payload is empty")
)

// Request inherits from http.Request, and provides additional methods.
type Request struct {
	*http.Request
}

// DecodeJSONPayload reads the request body and decodes the JSON using json.Unmarshal.
func (r *Request) DecodeJSONPayload(v interface{}) error {
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return ErrJSONPayloadEmpty
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}

// BaseURL returns a new URL object with the Host and Scheme taken from the request.
// (without the trailing slash in the host)
func (r *Request) BaseURL() *url.URL {
	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	host := r.Host
	if len(host) > 0 && host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

// URLFor returns the URL object from UriBase with the Path set to path, and the query
// string built with queryParams.
func (r *Request) URLFor(path string, queryParams map[string][]string) *url.URL {
	baseURL := r.BaseURL()
	baseURL.Path = path
	if queryParams != nil {
		query := url.Values{}
		for k, v := range queryParams {
			for _, vv := range v {
				query.Add(k, vv)
			}
		}
		baseURL.RawQuery = query.Encode()
	}
	return baseURL
}

// CorsInfo contains the CORS request info derived from a rest.Request.
type CorsInfo struct {
	IsCors      bool
	IsPreflight bool
	Origin      string
	OriginURL   *url.URL

	// The header value is converted to uppercase to avoid common mistakes.
	AccessControlRequestMethod string

	// The header values are normalized with http.CanonicalHeaderKey.
	AccessControlRequestHeaders []string
}

// GetCorsInfo derives CorsInfo from Request.
func (r *Request) GetCorsInfo() *CorsInfo {

	origin := r.Header.Get("Origin")

	var originURL *url.URL
	var isCors bool

	if origin == "" {
		isCors = false
	} else if origin == "null" {
		isCors = true
	} else {
		var err error
		originURL, err = url.ParseRequestURI(origin)
		isCors = err == nil && r.Host != originURL.Host
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")

	reqHeaders := []string{}
	rawReqHeaders := r.Header[http.CanonicalHeaderKey("Access-Control-Request-Headers")]
	for _, rawReqHeader := range rawReqHeaders {
		// net/http does not handle comma delimited headers for us
		for _, reqHeader := range strings.Split(rawReqHeader, ",") {
			reqHeaders = append(reqHeaders, http.CanonicalHeaderKey(strings.TrimSpace(reqHeader)))
		}
	}

	isPreflight := isCors && r.Method == "OPTIONS" && reqMethod != ""

	return &CorsInfo{
		IsCors:                      isCors,
		IsPreflight:                 isPreflight,
		Origin:                      origin,
		OriginURL:                   originURL,
		AccessControlRequestMethod:  strings.ToUpper(reqMethod),
		AccessControlRequestHeaders: reqHeaders,
	}
}
