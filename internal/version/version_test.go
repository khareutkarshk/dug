package version

import "testing"

func TestShortCommit(t *testing.T) {
	if got := shortCommit("abcdefghij"); got != "abcdefg" {
		t.Fatalf("got %q", got)
	}
	if got := shortCommit("abc"); got != "abc" {
		t.Fatalf("got %q", got)
	}
}

func TestInfoPrefersLdflags(t *testing.T) {
	origV, origC, origD := Version, Commit, Date
	t.Cleanup(func() {
		Version, Commit, Date = origV, origC, origD
	})

	Version = "9.9.9"
	Commit = "deadbee"
	Date = "2026-01-01T00:00:00Z"

	v, c, d := Info()
	if v != "9.9.9" || c != "deadbee" || d != "2026-01-01T00:00:00Z" {
		t.Fatalf("Info()=%s %s %s", v, c, d)
	}
}
