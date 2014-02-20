package wrapstesting

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-on/wrap"
)

type reverseProxy struct {
	*httputil.ReverseProxy
}

func ReverseProxy(rev *httputil.ReverseProxy) wrap.Wrapper {
	return &reverseProxy{rev}
}

func ReverseProxyByUrl(urlbase string) wrap.Wrapper {
	u, err := url.Parse(urlbase)

	if err != nil {
		panic(err.Error())
	}

	return &reverseProxy{httputil.NewSingleHostReverseProxy(u)}
}

// ServeHandle serves the request via the ReverseProxy, ignoring any inner handler
func (rp *reverseProxy) ServeHandle(inner http.Handler, rw http.ResponseWriter, req *http.Request) {
	rp.ServeHTTP(rw, req)
}

// Wrap wraps the given inner handler with the returned handler
func (rp *reverseProxy) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(rp, inner)
}
