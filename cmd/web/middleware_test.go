package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abdelrahman-habib/noter/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	rr, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeadersMiddleware(next).ServeHTTP(rr, r)

	rs := rr.Result()
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), "deny")
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), "0")
	assert.Equal(t, rs.Header.Get("Server"), "Go")
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(bytes.TrimSpace(body)), "OK")

}

func TestCacheControlMiddleware(t *testing.T) {
	rr, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	app := &application{}
	app.cacheControlMiddleware(next).ServeHTTP(rr, r)

	rs := rr.Result()
	assert.Equal(t, rs.Header.Get("Cache-Control"), "public, max-age=86400")
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(bytes.TrimSpace(body)), "OK")
}
