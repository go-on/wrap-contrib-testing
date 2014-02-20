package wrapstesting

import (
	"net/http"
	"testing"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
)

func anyway(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`anyway`))
}

type panicker struct{}

func (panicker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("don't panic")
}

func TestDefer(t *testing.T) {
	r := wrap.New(
		DeferFunc(anyway),
		wrap.Handler(panicker{}),
	)
	rw, req := helper.NewTestRequest("GET", "/")
	defer func() { recover() }()
	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw, "anyway", 200)
	if err != nil {
		t.Error(err.Error())
	}
}
