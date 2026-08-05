package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/khareutkarshk/dug/internal/version"
)

func TestPrintVersionDefault(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc1234"
	version.Date = "2026-08-05T00:00:00Z"
	version.Go = "go1.25.1"

	var buf bytes.Buffer
	if err := printVersion(&buf, nil); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	for _, want := range []string{
		"dug 1.2.3",
		"commit: abc1234",
		"built:  2026-08-05T00:00:00Z",
		"go:     go1.25.1",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output %q missing %q", out, want)
		}
	}
}

func TestPrintVersionShort(t *testing.T) {
	version.Version = "1.2.3"

	var buf bytes.Buffer
	if err := printVersion(&buf, []string{"-short"}); err != nil {
		t.Fatal(err)
	}

	if got := strings.TrimSpace(buf.String()); got != "1.2.3" {
		t.Fatalf("got %q, want %q", got, "1.2.3")
	}
}

func TestPrintVersionJSON(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc1234"
	version.Date = "2026-08-05T00:00:00Z"
	version.Go = "go1.25.1"

	var buf bytes.Buffer
	if err := printVersion(&buf, []string{"-json"}); err != nil {
		t.Fatal(err)
	}

	var info VersionInfo
	if err := json.Unmarshal(buf.Bytes(), &info); err != nil {
		t.Fatal(err)
	}

	if info.Version != "1.2.3" || info.Commit != "abc1234" {
		t.Fatalf("unexpected json payload: %+v", info)
	}
}
