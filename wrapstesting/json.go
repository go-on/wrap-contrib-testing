package wrapstesting

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-on/wrap"
)

// casts the Responsewriter to http.Handler in order to write to itself
type json_ struct {
	Type reflect.Type
}

func (j json_) newPtr() reflect.Value {
	val := reflect.New(j.Type)
	ref := reflect.New(val.Type())
	ref.Elem().Set(val)
	return ref
}

func (j json_) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.(http.Handler).ServeHTTP(w, r)
		ptr := j.newPtr()
		err := unWrap(w, ptr)
		if err != nil {
			panic(err.Error())
		}
		b, e := json.Marshal(ptr.Elem().Interface())
		if e != nil {
			panic(e.Error())
		}
		w.Write(b)
	})
}

func Json(t interface{}) wrap.Wrapper {
	ty := reflect.TypeOf(t)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return json_{ty}
}
