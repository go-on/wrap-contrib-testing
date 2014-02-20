package wrapstesting

import (
	"net/http"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
	. "launchpad.net/gocheck"
)

type deferSuite struct{}

var _ = Suite(&deferSuite{})

func anyway(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`anyway`))
}

type panicker struct{}

func (panicker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("don't panic")
}

func (s *deferSuite) TestDefer(c *C) {
	r := wrap.New(
		DeferFunc(anyway),
		wrap.Handler(panicker{}),
	)
	rw, req := helper.NewTestRequest("GET", "/")
	defer func() { recover() }()
	r.ServeHTTP(rw, req)
	helper.AssertResponse(rw, "anyway", 200)
}
