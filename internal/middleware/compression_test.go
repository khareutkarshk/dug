package middleware_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/middleware"
)

func TestCompressionGzipWhenAccepted(t *testing.T) {
	payload := strings.Repeat("abcdefghij", 200) // 2000 bytes

	h := middleware.Compression(config.CompressionConfig{
		Enabled: true,
		MinSize: 100,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(payload))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("Content-Encoding=%q", rec.Header().Get("Content-Encoding"))
	}

	gr, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = gr.Close() }()

	got, err := io.ReadAll(gr)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != payload {
		t.Fatalf("decompressed length=%d want=%d", len(got), len(payload))
	}
}

func TestCompressionSkipsWithoutAcceptEncoding(t *testing.T) {
	payload := strings.Repeat("x", 500)

	h := middleware.Compression(config.CompressionConfig{
		Enabled: true,
		MinSize: 10,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(payload))
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Header().Get("Content-Encoding") != "" {
		t.Fatal("expected no compression")
	}
	if rec.Body.String() != payload {
		t.Fatal("body mismatch")
	}
}

func TestCompressionSkipsSmallResponses(t *testing.T) {
	h := middleware.Compression(config.CompressionConfig{
		Enabled: true,
		MinSize: 1024,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("tiny"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "" {
		t.Fatal("expected no compression for small body")
	}
}

func TestCompressionSkipsAlreadyEncoded(t *testing.T) {
	h := middleware.Compression(config.CompressionConfig{
		Enabled: true,
		MinSize: 1,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "br")
		_, _ = w.Write(bytes.Repeat([]byte("z"), 100))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "br" {
		t.Fatalf("Content-Encoding=%q", rec.Header().Get("Content-Encoding"))
	}
}

func TestCompressionSkipsWebSocketUpgrade(t *testing.T) {
	h := middleware.Compression(config.CompressionConfig{
		Enabled: true,
		MinSize: 1,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(bytes.Repeat([]byte("w"), 200))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "" {
		t.Fatal("websocket responses must not be gzipped")
	}
}

func TestCompressionSkipsSSEAccept(t *testing.T) {
	h := middleware.Compression(config.CompressionConfig{
		Enabled: true,
		MinSize: 1,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(bytes.Repeat([]byte("s"), 200))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept", "text/event-stream")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "" {
		t.Fatal("SSE must not be gzipped")
	}
}
