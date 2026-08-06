package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
)

func TestInitCreatesProject(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "my-gateway")

	var stdout, stderr bytes.Buffer
	if err := initProject(&stdout, &stderr, []string{target}); err != nil {
		t.Fatalf("init failed: %v\nstderr=%s", err, stderr.String())
	}

	expected := []string{
		filepath.Join("configs", "edge.yaml"),
		"docker-compose.yml",
		".gitignore",
		filepath.Join("certs", ".gitkeep"),
		"README.md",
	}

	for _, rel := range expected {
		path := filepath.Join(target, rel)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}

	cfg, err := config.Load(filepath.Join(target, "configs", "edge.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 8080 || len(cfg.Routes) != 1 {
		t.Fatalf("unexpected config: %+v", cfg)
	}

	compose, err := os.ReadFile(filepath.Join(target, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(compose)
	if !strings.Contains(text, defaultGHCRImage) {
		t.Fatalf("compose missing GHCR image:\n%s", text)
	}
	if !strings.Contains(text, `container_name: "my-gateway"`) {
		t.Fatalf("compose missing project name:\n%s", text)
	}

	out := stdout.String()
	for _, want := range []string{
		`Created DUG project "my-gateway"`,
		"cd " + target,
		"dug run -config configs/edge.yaml",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
}

func TestInitDoesNotOverwriteWithoutForce(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "gw")

	var stdout, stderr bytes.Buffer
	if err := initProject(&stdout, &stderr, []string{target}); err != nil {
		t.Fatal(err)
	}

	marker := filepath.Join(target, "configs", "edge.yaml")
	if err := os.WriteFile(marker, []byte("custom: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	err := initProject(&stdout, &stderr, []string{target})
	if err == nil {
		t.Fatal("expected error when all files exist")
	}

	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "custom: true\n" {
		t.Fatalf("file was overwritten without --force: %q", data)
	}
}

func TestInitForceOverwrites(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "gw")

	var stdout, stderr bytes.Buffer
	if err := initProject(&stdout, &stderr, []string{target}); err != nil {
		t.Fatal(err)
	}

	marker := filepath.Join(target, "README.md")
	if err := os.WriteFile(marker, []byte("stale\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	if err := initProject(&stdout, &stderr, []string{target, "--force"}); err != nil {
		t.Fatalf("force init failed: %v\nstderr=%s", err, stderr.String())
	}

	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(data)) == "stale" {
		t.Fatal("expected README to be overwritten with --force")
	}
	if !strings.Contains(string(data), "# gw") {
		t.Fatalf("unexpected README content: %s", data)
	}
}

func TestInitRequiresDirectory(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := initProject(&stdout, &stderr, nil)
	if err == nil {
		t.Fatal("expected error for missing directory")
	}
}
