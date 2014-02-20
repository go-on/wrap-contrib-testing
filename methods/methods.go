package methods

import (
	"fmt"
	"net/http"
)

//type Method string
//type verb int

type Method int

const (
	POST Method = 1 << iota
	GET
	PUT
	DELETE
	PATCH
	OPTIONS
	HEAD
	TRACE
)

/*
const (
	POST    Method = "POST"
	GET     Method = "GET"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	OPTIONS Method = "OPTIONS"
	HEAD    Method = "HEAD"
	TRACE   Method = "TRACE"
)
*/

var methodStrings = map[Method]string{
	POST:    "POST",
	GET:     "GET",
	PUT:     "PUT",
	DELETE:  "DELETE",
	PATCH:   "PATCH",
	OPTIONS: "OPTIONS",
	HEAD:    "HEAD",
	TRACE:   "TRACE",
}

var StringToMethod = map[string]Method{
	"POST":    POST,
	"GET":     GET,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"PATCH":   PATCH,
	"OPTIONS": OPTIONS,
	"HEAD":    HEAD,
	"TRACE":   TRACE,
}

func (m Method) String() string {
	return methodStrings[m]
}

/*
The following definitions are based on
http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html

and for PATCH
http://tools.ietf.org/html/rfc5789
*/

// a method is considered safe, if it
// has no requested sideeffect by the user agent
func (m Method) IsSafe() bool {
	return m == GET || m == HEAD || m == OPTIONS || m == TRACE
}

// a method is idempotent, if (aside from error or expiration issues)
// the side-effects of N > 0 identical requests is the same as for a single request
func (m Method) IsIdempotent() bool {
	return m != POST && m != PATCH
}

func (m Method) IsResponseCacheable() bool {
	// if it meets the requirements for HTTP caching described in section 13
	return m == GET || m == HEAD
}

// the the method return an empty message body
func (m Method) EmptyBody() bool {
	return m == HEAD || m == OPTIONS
}

// RequestMethod is a shortcut to get the method of a *http.Request
// it panicks if the method is unknown. if you want to check, if a method
// is known, use StringToMethod
func RequestMethod(rq *http.Request) Method {
	m, has := StringToMethod[rq.Method]
	if !has {
		panic(fmt.Sprintf("unknown HTTP method %#v", rq.Method))
	}
	return m
}

/*
if the new field values indicate that the cached entity differs from the current entity (as would be indicated by a change in Content-Length, Content-MD5, ETag or Last-Modified), then the cache MUST treat the cache entry as stale.
*/

/*
The semantics of the GET method change to a "conditional GET" if the request message includes an If-Modified-Since, If-Unmodified-Since, If-Match, If-None-Match, or If-Range header field. A conditional GET method requests that the entity be transferred only under the circumstances described by the conditional header field(s). The conditional GET method is intended to reduce unnecessary network usage by allowing cached entities to be refreshed without requiring multiple requests or transferring data already held by the client.


The semantics of the GET method change to a "partial GET" if the request message includes a Range header field. A partial GET requests that only part of the entity be transferred, as described in section 14.35. The partial GET method is intended to reduce unnecessary network usage by allowing partially-retrieved entities to be completed without transferring data already held by the client.
*/

/*
The metainformation contained in the HTTP headers in response to a HEAD request SHOULD be identical to the information sent in response to a GET request
*/
