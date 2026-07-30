package cli

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
)

func checkUpstreams(cfg *config.Config) {

	fmt.Println("Upstreams")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, route := range cfg.Routes {

		for _, upstream := range route.Upstreams {

			checkUpstream(client, upstream.URL)
		}
	}
}

func checkUpstream(client *http.Client, rawURL string) {

	u, err := url.Parse(rawURL)
	if err != nil {

		fmt.Printf("✗ %s invalid URL\n", rawURL)
		return
	}

	fmt.Printf("• %s\n", rawURL)

	resp, err := client.Get(rawURL)

	if err != nil {

		fmt.Println("  ✗ unreachable")
		return
	}

	_ = resp.Body.Close()

	fmt.Printf("  ✓ reachable (%d)\n", resp.StatusCode)

	healthURL := *u
	healthURL.Path = "/health"

	resp, err = client.Get(healthURL.String())

	if err != nil {

		fmt.Println("  ✗ health endpoint unreachable")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		fmt.Println("  ✓ healthy")
	} else {

		fmt.Printf("  ✗ unhealthy (%d)\n", resp.StatusCode)
	}
}
