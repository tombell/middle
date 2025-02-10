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

func TestRecovery(t *testing.T) {
	is := is.New(t)

	buf := bytes.NewBufferString("")
	logger := slog.New(slog.NewTextHandler(buf, nil))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom town")
	})

	fn := middle.Use(
		ware.Logger(logger),
		ware.Recovery(),
	)(handler)

	req := httptest.NewRequest("GET", "https://example.com/", nil)
	resp := httptest.NewRecorder()

	fn.ServeHTTP(resp, req)

	is.Equal(resp.Code, http.StatusInternalServerError)
	is.True(strings.Contains(buf.String(), `level=ERROR msg="recovered from panic" err="boom town"`+"\n"))
}
