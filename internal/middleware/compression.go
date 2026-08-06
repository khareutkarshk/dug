package middleware

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"net"
	"net/http"
	"strings"

	"github.com/khareutkarshk/dug/internal/config"
)

// Compression optionally gzip-encodes responses when the client accepts gzip.
func Compression(cfg config.CompressionConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled || !acceptsGzip(r) || skipCompressionRequest(r) {
				next.ServeHTTP(w, r)
				return
			}

			minSize := cfg.MinSize
			if minSize < 0 {
				minSize = 0
			}

			cw := &compressWriter{
				ResponseWriter: w,
				minSize:        minSize,
				status:         http.StatusOK,
			}
			defer cw.close()

			next.ServeHTTP(cw, r)
		})
	}
}

func acceptsGzip(r *http.Request) bool {
	for _, part := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
		encoding := strings.TrimSpace(strings.Split(part, ";")[0])
		if strings.EqualFold(encoding, "gzip") {
			return true
		}
	}
	return false
}

func skipCompressionRequest(r *http.Request) bool {
	if isWebSocketRequest(r) {
		return true
	}
	accept := r.Header.Get("Accept")
	return strings.Contains(strings.ToLower(accept), "text/event-stream")
}

func isWebSocketRequest(r *http.Request) bool {
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return false
	}
	for _, part := range strings.Split(r.Header.Get("Connection"), ",") {
		if strings.EqualFold(strings.TrimSpace(part), "upgrade") {
			return true
		}
	}
	return false
}

type compressWriter struct {
	http.ResponseWriter

	minSize int
	status  int

	buf           bytes.Buffer
	gz            *gzip.Writer
	headerWritten bool
	passthrough   bool
	gzipEnabled   bool
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	cw.status = statusCode
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	if cw.passthrough {
		if err := cw.ensureHeader(); err != nil {
			return 0, err
		}
		return cw.ResponseWriter.Write(p)
	}

	if cw.shouldDisable() {
		cw.passthrough = true
		if err := cw.flushBufferPlain(); err != nil {
			return 0, err
		}
		return cw.ResponseWriter.Write(p)
	}

	if cw.gzipEnabled {
		return cw.gz.Write(p)
	}

	if _, err := cw.buf.Write(p); err != nil {
		return 0, err
	}

	if cw.buf.Len() >= cw.minSize {
		if err := cw.enableGzip(); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func (cw *compressWriter) Flush() {
	if cw.gzipEnabled {
		if cw.gz != nil {
			_ = cw.gz.Flush()
		}
		if f, ok := cw.ResponseWriter.(http.Flusher); ok {
			f.Flush()
		}
		return
	}

	if cw.passthrough {
		if f, ok := cw.ResponseWriter.(http.Flusher); ok {
			f.Flush()
		}
		return
	}

	// Ignore empty flushes (ReverseProxy may flush before body writes).
	if cw.buf.Len() == 0 {
		return
	}

	// A flush with buffered data means streaming — skip compression.
	cw.passthrough = true
	_ = cw.flushBufferPlain()
	if f, ok := cw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (cw *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := cw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

func (cw *compressWriter) close() {
	if cw.passthrough {
		return
	}
	if cw.gzipEnabled {
		if cw.gz != nil {
			_ = cw.gz.Close()
		}
		return
	}
	_ = cw.flushBufferPlain()
}

func (cw *compressWriter) shouldDisable() bool {
	if enc := cw.Header().Get("Content-Encoding"); enc != "" {
		return true
	}
	ct := strings.ToLower(cw.Header().Get("Content-Type"))
	return strings.Contains(ct, "text/event-stream")
}

func (cw *compressWriter) enableGzip() error {
	cw.Header().Del("Content-Length")
	cw.Header().Set("Content-Encoding", "gzip")
	cw.Header().Add("Vary", "Accept-Encoding")

	if err := cw.ensureHeader(); err != nil {
		return err
	}

	cw.gz = gzip.NewWriter(cw.ResponseWriter)
	cw.gzipEnabled = true

	if cw.buf.Len() > 0 {
		if _, err := cw.gz.Write(cw.buf.Bytes()); err != nil {
			return err
		}
		cw.buf.Reset()
	}
	return nil
}

func (cw *compressWriter) flushBufferPlain() error {
	if err := cw.ensureHeader(); err != nil {
		return err
	}
	if cw.buf.Len() == 0 {
		return nil
	}
	_, err := cw.ResponseWriter.Write(cw.buf.Bytes())
	cw.buf.Reset()
	return err
}

func (cw *compressWriter) ensureHeader() error {
	if cw.headerWritten {
		return nil
	}
	cw.ResponseWriter.WriteHeader(cw.status)
	cw.headerWritten = true
	return nil
}
