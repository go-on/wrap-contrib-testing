package wrapstesting

import (
	"net/http"

	"github.com/go-on/wrap"
)

type (
	// defer something after the request handling
	defer_ struct{ http.Handler }
)

func (ø defer_) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { ø.ServeHTTP(w, r) }()
		in.ServeHTTP(w, r)
	})
}

func DeferFunc(fn func(http.ResponseWriter, *http.Request)) wrap.Wrapper {
	return defer_{http.HandlerFunc(fn)}
}

func Defer(h http.Handler) wrap.Wrapper {
	return defer_{h}
}
