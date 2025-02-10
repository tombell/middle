package ware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/tombell/middle"
	"github.com/tombell/middle/ware"
)

func TestRequestLogging(t *testing.T) {
	is := is.New(t)

	buf := bytes.NewBufferString("")
	logger := slog.New(slog.NewTextHandler(buf, nil))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	})

	fn := middle.Use(
		ware.Logger(logger),
		ware.RequestLogging(),
	)(handler)

	req := httptest.NewRequest("GET", "https://example.com/", nil)
	resp := httptest.NewRecorder()

	fn.ServeHTTP(resp, req)

	is.True(strings.Contains(buf.String(), `level=INFO msg=http:started method=GET path=/`+"\n"))
	is.True(strings.Contains(buf.String(), `level=INFO msg=http:finished method=GET path=/ status=201 size=7 elapsed=`))
}
