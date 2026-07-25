package router

import (
	"net/http"
	"sync/atomic"
)

type Manager struct {
	handler atomic.Value
}

func NewManager(h http.Handler) *Manager {

	m := &Manager{}
	m.handler.Store(h)

	return m
}

func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := m.handler.Load().(http.Handler)
	h.ServeHTTP(w, r)
}

func (m *Manager) Update(h http.Handler) {
	m.handler.Store(h)
}
