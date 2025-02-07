package middle_test

import (
	"net/http"
	"testing"

	"github.com/matryer/is"

	"github.com/tombell/middle"
)

func TestMiddleUse(t *testing.T) {
	result := ""

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result += "handler"
	})

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "one "
			next.ServeHTTP(w, r)
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "two "
			next.ServeHTTP(w, r)
		})
	}

	fn := middle.Use(mw1, mw2)
	fn(handle).ServeHTTP(nil, nil)

	is := is.New(t)
	is.Equal("one two handler", result) // expect to be the same
}
