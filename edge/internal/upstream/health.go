package upstream

import (
	"net/http"
	"time"
)

func (p *Pool) StartHealthCheck(interval time.Duration) {

	go func() {

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			for _, backend := range p.backends {
				p.checkBackend(backend)
			}
		}
	}()
}

func (p *Pool) checkBackend(backend *Backend) {

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(backend.URL.String() + "/health")
	if err != nil {
		backend.ReportFailure()
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		backend.ReportSuccess()
		backend.Healthy.Store(true)
	} else {
		backend.ReportFailure()
	}
}
