package proxy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/upstream"
)

type Proxy struct {
	proxy *httputil.ReverseProxy
}

func New(pool *upstream.Pool, retries int, timeout time.Duration, requestHeaders config.HeaderRules, responseHeaders config.HeaderRules) *Proxy {

	transport := &RetryTransport{
		Base: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,
		},
		Pool:            pool,
		Retries:         retries,
		Timeout:         timeout,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: responseHeaders,
	}

	rp := &httputil.ReverseProxy{

		// RetryTransport will choose one for every attempt.
		Director: func(req *http.Request) {},

		Transport: transport,

		ModifyResponse: func(resp *http.Response) error {

			// add response headers
			for k, v := range transport.ResponseHeaders.Add {
				resp.Header.Set(k, v)
			}

			// remove response headers
			for _, h := range transport.ResponseHeaders.Remove {
				resp.Header.Del(h)
			}

			return nil
		},

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {

			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
				return
			}

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
