package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/khareutkarshk/dug/edge/internal/upstream"
)

func New(pool *upstream.Pool) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			backend := pool.Next()

			if backend == nil {
				return
			}

			// Rewrite the request so it is forwarded
			// to the selected backend.
			req.URL.Scheme = backend.URL.Scheme
			req.URL.Host = backend.URL.Host

			// Update the Host header as well.
			req.Host = backend.URL.Host
		},
	}
}
