package wrapstesting

import (
	"net/http"

	"github.com/go-on/rack"
)

type http1_1 struct{}

var HTTP1_1 = http1_1{}

func (h http1_1) ServeHandle(in http.Handler, w http.ResponseWriter, r *http.Request) {
	if !r.ProtoAtLeast(1, 1) {
		// protocol not supported
		w.WriteHeader(505)
		return
	}
	in.ServeHTTP(w, r)
}

func (h http1_1) Wrap(in http.Handler) http.Handler {
	return rack.ServeHandle(h, in)
}
