package testutil

import (
	"net/http"
	"net/http/httptest"
)

func NewBackend(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}
