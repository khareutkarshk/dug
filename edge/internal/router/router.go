package router

import (
	"net/http"

	"github.com/khareutkarshk/dug/edge/internal/proxy"
)

func NewRouter() (http.Handler, error) {
	mux := http.NewServeMux()

	p, err := proxy.New("http://localhost:3001")

	if err != nil {
		return nil, err
	}

	mux.Handle("/", p)
	return mux, nil
}
