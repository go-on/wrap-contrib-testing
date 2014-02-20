package wrapstesting

import (
	"net/http"

	"github.com/go-on/wrap"
)

// casts the Responsewriter to http.Handler in order to write to itself
type responseWriterHandle struct{}

func (ø responseWriterHandle) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.(http.Handler).ServeHTTP(w, r)
	})
}

var ResponseWriterHandler wrap.Wrapper = responseWriterHandle{}
