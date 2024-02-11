package ware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/tombell/middle/ware"
)

func TestLogger(t *testing.T) {
	is := is.New(t)

	buf := bytes.NewBufferString("")
	logger := slog.New(slog.NewTextHandler(buf, nil))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, _ := r.Context().Value(ware.LoggerContextKey).(*slog.Logger)
		is.Equal(logger, log)
		log.Info("hello world")
	})

	req := httptest.NewRequest("GET", "https://example.com", nil)
	resp := httptest.NewRecorder()

	fn := ware.Logger(logger)(handler)
	fn.ServeHTTP(resp, req)

	is.True(strings.Contains(buf.String(), "level=INFO msg=\"hello world\""))
}
