package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// client stores a limiter and the last time it was used.
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Manager manages one rate limiter per client IP.
type Manager struct {
	mu sync.RWMutex

	rps   rate.Limit
	burst int

	clients map[string]*client
}

// NewManager creates a new rate limiter manager.
func NewManager(rps float64, burst int) *Manager {

	manager := &Manager{
		rps:     rate.Limit(rps),
		burst:   burst,
		clients: make(map[string]*client),
	}

	// Start background cleanup of idle clients.
	go manager.cleanup()

	return manager
}

// Get returns the limiter for a client.
// A new limiter is created on the first request.
func (m *Manager) Get(ip string) *rate.Limiter {

	// Fast path: check if limiter already exists.
	m.mu.RLock()
	entry, ok := m.clients[ip]
	m.mu.RUnlock()

	if ok {
		// Update last activity time.
		m.mu.Lock()
		entry.lastSeen = time.Now()
		m.mu.Unlock()

		return entry.limiter
	}

	// Slow path: create a new limiter.
	m.mu.Lock()
	defer m.mu.Unlock()

	// Another goroutine may have created it while we waited.
	if entry, ok := m.clients[ip]; ok {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(m.rps, m.burst)

	m.clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

// cleanup periodically removes idle clients to free memory.
func (m *Manager) cleanup() {

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {

		m.mu.Lock()

		for ip, client := range m.clients {

			// Remove clients idle for more than 10 minutes.
			if time.Since(client.lastSeen) > 10*time.Minute {
				delete(m.clients, ip)
			}
		}

		m.mu.Unlock()
	}
}
