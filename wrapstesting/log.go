package wrapstesting

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-on/rack"
	"github.com/go-on/rack/helper"
	"github.com/go-on/wrap"
)

type logger struct{ *log.Logger }

func (l *logger) ServeHandle(in http.Handler, w http.ResponseWriter, r *http.Request) {
	// l.Logger.Printf("ResponseWriter: %#v\nRequest: %#v\n", w, r)
	requestHeaders := fmt.Sprintf("%v", r.Header)

	fake := helper.NewFake(w)

	in.ServeHTTP(fake, r)

	if fake.HasChanged() {
		l.Printf(`
-- REQUEST --
%s %s
HEADERS
%s
-- RESPONSE --
STATUS CODE: %d
HEADERS
%s
BODY
%s

`, r.Method, r.URL.Path, requestHeaders, fake.WHeader, fake.Header(), string(fake.Buffer.Bytes()))
	}

	fake.WriteHeaderTo(w)
	if fake.WHeader != 0 {
		w.WriteHeader(fake.WHeader)
	}

	fake.Buffer.WriteTo(w)
}

func (l *logger) Wrap(inner http.Handler) http.Handler {
	return rack.ServeHandle(l, inner)
}

func LOGGER(prefix string) wrap.Wrapper {
	return Logger(log.New(os.Stderr, prefix+" ", log.LstdFlags))
}

func Logger(l *log.Logger) wrap.Wrapper {
	return &logger{l}
}
