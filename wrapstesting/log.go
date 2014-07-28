package wrapstesting

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-on/wrap"
)

type logger struct{ *log.Logger }

func (l *logger) ServeHandle(in http.Handler, w http.ResponseWriter, r *http.Request) {
	checked := wrap.NewRWPeek(w, func(ck *wrap.RWPeek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})
	// l.Logger.Printf("ResponseWriter: %#v\nRequest: %#v\n", w, r)
	requestHeaders := fmt.Sprintf("%v", r.Header)

	// buf := helper.NewResponseBuffer(w)

	in.ServeHTTP(checked, r)

	if checked.HasChanged() {
		l.Printf(`
-- REQUEST --
%s %s
HEADERS
%s
-- RESPONSE --
STATUS CODE: %d
HEADERS
%s
`, r.Method, r.URL.Path, requestHeaders, checked.Code, checked.Header())
	}

}

func (l *logger) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(l, inner)
}

func LOGGER(prefix string) wrap.Wrapper {
	return Logger(log.New(os.Stderr, prefix+" ", log.LstdFlags))
}

func Logger(l *log.Logger) wrap.Wrapper {
	return &logger{l}
}
