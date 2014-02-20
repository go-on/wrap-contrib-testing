package wrapstesting

import (
	"fmt"
	"testing"

	// "fmt"
	"net/http"
	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type restSuite struct{}

var _ = Suite(&restSuite{})

func methodWrite(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, req.Method)
}

func (s *restSuite) TestRest(c *C) {
	r := wrap.New(
		MethodOverride,
		wrap.HandlerFunc(methodWrite),
	)

	for overrideMethod, requestMethod := range acceptedOverrides {
		rw, req := helper.NewTestRequest(requestMethod, "/")
		req.Header.Set(overrideHeader, overrideMethod)
		r.ServeHTTP(rw, req)
		err := helper.AssertResponse(rw, overrideMethod, 200)
		c.Assert(err, Equals, nil)
		h := req.Header.Get(overrideHeader)
		c.Assert(h, Equals, "")
	}
}
