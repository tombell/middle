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

func TestRequestID(t *testing.T) {
	is := is.New(t)

	buf := bytes.NewBufferString("")
	logger := slog.New(slog.NewTextHandler(buf, nil))

	result := ""

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, _ := r.Context().Value(ware.LoggerContextKey).(*slog.Logger)
		log.Info("hello world")
		result = r.Context().Value(ware.RequestIDContextKey).(string)
	})

	fn := middle.Use(
		ware.Logger(logger),
		ware.RequestID(func() string { return "96ad9a25-efcb-4b84-b7cd-c09cc166686a" }),
	)(handler)

	req := httptest.NewRequest("GET", "https://example.com", nil)
	resp := httptest.NewRecorder()

	fn.ServeHTTP(resp, req)

	is.Equal("96ad9a25-efcb-4b84-b7cd-c09cc166686a", result)
	is.Equal("96ad9a25-efcb-4b84-b7cd-c09cc166686a", resp.Header().Get("Request-ID"))
	is.True(strings.Contains(buf.String(), `level=INFO msg="hello world" rid=96ad9a25-efcb-4b84-b7cd-c09cc166686a`))
}
