package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	s := String()
	if !strings.Contains(s, "taskcapsule") {
		t.Errorf("expected taskcapsule in version string, got %q", s)
	}
	if !strings.Contains(s, Version) {
		t.Errorf("expected version %q in string, got %q", Version, s)
	}
}

func TestVersionConstants(t *testing.T) {
	if Version == "" {
		t.Error("Version must not be empty")
	}
}
