package proxy

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/khareutkarshk/dug/internal/upstream"
)

type Proxy struct {
	proxy *httputil.ReverseProxy
}

func New(pool *upstream.Pool, retries int) *Proxy {

	transport := &RetryTransport{
		Base: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,
		},
		Pool:    pool,
		Retries: retries,
	}

	rp := &httputil.ReverseProxy{

		// RetryTransport will choose one for every attempt.
		Director: func(req *http.Request) {},

		Transport: transport,

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		},
	}

	return &Proxy{
		proxy: rp,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
