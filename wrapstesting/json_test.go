package wrapstesting

import (
	// "fmt"
	"net/http"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
	"github.com/go-on/wrap-contrib/wraps"
	. "launchpad.net/gocheck"
)

type jsonSuite struct{}

var _ = Suite(&jsonSuite{})

func mkJsonCtx(rw http.ResponseWriter, req *http.Request) http.ResponseWriter {
	return &jsonCtx{ResponseWriter: rw}
}

type jsonCtx struct {
	Path                string
	http.ResponseWriter `json:"-"`
}

func (j *jsonCtx) Prepare(w http.ResponseWriter, r *http.Request) {
	j.Path = r.URL.Path
}

func (s *jsonSuite) TestJson(c *C) {

	r := wrap.New(
		Context(mkJsonCtx),
		wraps.Before(HandlerMethod((*jsonCtx).Prepare)),
		Json(&jsonCtx{}),
	)

	rw, req := helper.NewTestRequest("GET", "/hiho")
	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw, `{"Path":"/hiho"}`, 200)

	c.Assert(err, Equals, nil)

	rw, req = helper.NewTestRequest("GET", "/huho")
	r.ServeHTTP(rw, req)
	err = helper.AssertResponse(rw, `{"Path":"/huho"}`, 200)

	c.Assert(err, Equals, nil)
}
