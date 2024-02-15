package ware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"

	"github.com/tombell/middle/ware"
)

func TestCORS(t *testing.T) {
	is := is.New(t)

	opts := ware.CORSOptions{
		AllowedOrigins:   []string{"example.com", "localhost"},
		AllowedHeaders:   []string{"Header1", "Header2", "Header3"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowCredentials: true,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	fn := ware.CORS(opts)

	req := httptest.NewRequest("OPTIONS", "https://example.com", nil)
	req.Header.Add("Origin", "example.com")
	resp := httptest.NewRecorder()

	fn(handler).ServeHTTP(resp, req)

	is.Equal(http.StatusOK, resp.Code)
	is.Equal("example.com", resp.Header().Get("Access-Control-Allow-Origin"))
	is.Equal("Header1, Header2, Header3", resp.Header().Get("Access-Control-Allow-Headers"))
	is.Equal("GET, POST, PUT, DELETE, HEAD, OPTIONS", resp.Header().Get("Access-Control-Allow-Methods"))
	is.Equal("true", resp.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORSWildcard(t *testing.T) {
	is := is.New(t)

	opts := ware.CORSOptions{
		AllowedHeaders: []string{"Header1", "Header2", "Header3"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	fn := ware.CORS(opts)

	req := httptest.NewRequest("OPTIONS", "https://example.com", nil)
	req.Header.Add("Origin", "localhost")
	resp := httptest.NewRecorder()

	fn(handler).ServeHTTP(resp, req)

	is.Equal(http.StatusOK, resp.Code)
	is.Equal("*", resp.Header().Get("Access-Control-Allow-Origin"))
	is.Equal("Header1, Header2, Header3", resp.Header().Get("Access-Control-Allow-Headers"))
	is.Equal("GET, POST, PUT, DELETE, HEAD, OPTIONS", resp.Header().Get("Access-Control-Allow-Methods"))
	is.Equal("", resp.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORSNotAllowed(t *testing.T) {
	is := is.New(t)

	opts := ware.CORSOptions{
		AllowedOrigins: []string{"myhost"},
		AllowedHeaders: []string{"Header1", "Header2", "Header3"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	fn := ware.CORS(opts)

	req := httptest.NewRequest("OPTIONS", "https://example.com", nil)
	req.Header.Add("Origin", "localhost")
	resp := httptest.NewRecorder()

	fn(handler).ServeHTTP(resp, req)

	is.Equal(http.StatusOK, resp.Code)
	is.Equal("", resp.Header().Get("Access-Control-Allow-Origin"))
	is.Equal("Header1, Header2, Header3", resp.Header().Get("Access-Control-Allow-Headers"))
	is.Equal("GET, POST, PUT, DELETE, HEAD, OPTIONS", resp.Header().Get("Access-Control-Allow-Methods"))
	is.Equal("", resp.Header().Get("Access-Control-Allow-Credentials"))
}
