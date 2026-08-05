package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
)

func TestCheckConfigValid(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Routes: []config.Route{{
			Path:      "/",
			Strategy:  "least_connections",
			Upstreams: []config.Upstream{{URL: "http://localhost:3001", Weight: 1}},
		}},
	}

	result := checkConfig(cfg)
	if !result.Ok() {
		t.Fatalf("expected valid config, got: %v", result.Error())
	}
}

func TestCheckConfigCollectsIssues(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 0},
		Routes: []config.Route{
			{Path: "api", Strategy: "round_robin", Upstreams: nil},
			{Path: "/ok", Upstreams: []config.Upstream{{URL: "ftp://bad", Weight: -1}}},
		},
	}

	result := checkConfig(cfg)
	if result.Ok() {
		t.Fatal("expected validation issues")
	}

	joined := result.Error().Error()
	for _, want := range []string{
		"server.port",
		"routes[0].path",
		"routes[0].strategy",
		"routes[0].upstreams",
		"routes[1].upstreams[0].url",
		"routes[1].upstreams[0].weight",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing issue for %s in:\n%s", want, joined)
		}
	}
}

func TestValidateCommandJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ok.yaml")
	content := []byte(`
server:
  port: 8080
routes:
  - path: /
    upstreams:
      - url: http://127.0.0.1:9
        weight: 1
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := validate(&stdout, &stderr, []string{"-config", path, "-json"}); err != nil {
		t.Fatal(err)
	}

	var rep validateReport
	if err := json.Unmarshal(stdout.Bytes(), &rep); err != nil {
		t.Fatal(err)
	}
	if !rep.Valid || rep.Routes != 1 || rep.Server.Port != 8080 {
		t.Fatalf("unexpected report: %+v", rep)
	}
}

func TestValidateCommandInvalidQuiet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	content := []byte(`
server:
  port: 8080
routes: []
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := validate(&stdout, &stderr, []string{"-config", path, "-quiet"})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(stderr.String(), "invalid") {
		t.Fatalf("expected invalid message, got %q", stderr.String())
	}
}

func TestDoctorUnreachableUpstreamFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doctor.yaml")
	content := []byte(`
server:
  port: 18080
routes:
  - path: /
    upstreams:
      - url: http://127.0.0.1:1
        weight: 1
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := doctor(&stdout, &stderr, []string{
		"-config", path,
		"-timeout", (50 * time.Millisecond).String(),
		"-json",
	})
	if err == nil {
		t.Fatal("expected doctor failure for unreachable upstream")
	}

	var report DoctorReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Failed == 0 {
		t.Fatalf("expected failing checks, got %+v", report)
	}
}
