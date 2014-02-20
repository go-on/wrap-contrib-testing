package wrapstesting

import (
	"fmt"
	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib-testing/methods"
	"github.com/go-on/wrap-contrib/helper"

	"net/http"
)

// does something before the request is handled further
type methodOverride struct{}

var acceptedOverrides = map[string]string{
	"PATCH":   "POST",
	"OPTIONS": "GET",
	"DELETE":  "POST",
	"PUT":     "POST",
}

var overrideHeader = "X-HTTP-Method-Override"

func (ø methodOverride) ServeHandle(in http.Handler, w http.ResponseWriter, r *http.Request) {
	override := r.Header.Get(overrideHeader)

	if override != "" {
		expectedMethod, accepted := acceptedOverrides[override]
		if !accepted {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `Unsupported value for %s: %#v.
Supported values are PUT, DELETE, PATCH and OPTIONS`, overrideHeader, override)
			return
		}

		if expectedMethod != r.Method {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `%s with value %s only allowed for %s requests.`,
				overrideHeader, override, expectedMethod)
			return
		}

		// fmt.Printf("override method to %s\n", override)

		// everything went fine, override the method
		r.Header.Del(overrideHeader)
		r.Method = override
	}

	in.ServeHTTP(w, r)

}

func (ø methodOverride) Wrap(in http.Handler) (out http.Handler) {
	return wrap.ServeHandle(ø, in)
}

var MethodOverride = methodOverride{}

type filterBody struct {
	methods methods.Method
}

func (f *filterBody) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ver, _ := methods.StringToMethod[r.Method]
		if f.methods&ver == 0 {
			in.ServeHTTP(w, r)
			return
		}

		buf := helper.NewResponseBuffer()
		in.ServeHTTP(buf, r)
		buf.WriteHeadersTo(w)

		if buf.Code != 0 {
			w.WriteHeader(buf.Code)
		}
	})
}

// Filter the body for the given method(s)
// to filter mutiple methods, use FilterBody(methods.PATCH|methods.OPTIONS)
func FilterBody(m methods.Method) wrap.Wrapper {
	return &filterBody{m}
}

/*
type options {

}
*/

// Allow: HEAD,GET,PUT,DELETE,OPTIONS

/*
TODO: ignore body for options
*/
