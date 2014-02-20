package wrapstesting

import (
	"fmt"
	"net/http"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
	"github.com/go-on/wrap-contrib/wraps"
	. "launchpad.net/gocheck"
)

type contextSuite struct{}

var _ = Suite(&contextSuite{})

type ctx struct {
	path string
	http.ResponseWriter
}

type blindctx struct {
	http.ResponseWriter
}

func (c *ctx) Prepare(w http.ResponseWriter, r *http.Request) {
	c.path = r.URL.Path
}

func (c *ctx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("~" + c.path + "~"))
}

func check(w http.ResponseWriter, r *http.Request) {
	c := &ctx{}
	MustUnWrap(w, &c)
	w.Write([]byte("#" + c.path + "#"))
}

func (s *contextSuite) TestContextHandlerMethod(c *C) {
	r := wrap.New(
		Context(&blindctx{}),
		Context(&ctx{}),
		wraps.After(http.HandlerFunc(check)),
		wraps.Before(HandlerMethod((*ctx).Prepare)),
		wrap.Handler(HandlerMethod((*ctx).ServeHTTP)),
	)

	rw, req := helper.NewTestRequest("GET", "/path")
	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw, "~/path~#/path#", 200)

	c.Assert(err, Equals, nil)
}

func (s *contextSuite) TestContextUnwrapIdentical(c *C) {
	c1 := ctx{path: "x"}
	c2 := &ctx{}

	err := UnWrap(&c1, &c2)
	c.Assert(c2.path, Equals, "x")
	if err != nil {
		panic(err.Error())
	}

	c3 := &ctx{}
	err = UnWrap(c1, &c3)
	c.Assert(c3.path, Equals, "x")

	if err != nil {
		panic(err.Error())
	}
}

func (s *contextSuite) TestContextUnwrapNested(c *C) {
	c1 := blindctx{&ctx{path: "x"}}
	c2 := &ctx{}

	err := UnWrap(&c1, &c2)
	if err != nil {
		panic(err.Error())
	}
	c.Assert(c2.path, Equals, "x")

	c3 := &ctx{}

	err = UnWrap(c1, &c3)
	if err != nil {
		panic(err.Error())
	}
	c.Assert(c3.path, Equals, "x")
}

func (s *contextSuite) TestContextUnwrapError(c *C) {
	_ = fmt.Println
	c1 := blindctx{}
	c2 := &ctx{}
	err := UnWrap(&c1, &c2)
	// fmt.Println(err.Error())
	c.Assert(err, NotNil)

	rw, _ := helper.NewTestRequest("GET", "/path")
	c1 = blindctx{rw}
	err = UnWrap(&c1, &c2)
	// fmt.Println(err.Error())
	c.Assert(err, NotNil)
}
