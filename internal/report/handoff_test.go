package report

import (
	"strings"
	"testing"
)

func TestGenerateHandoff(t *testing.T) {
	data := HandoffData{
		Name:            "payment-timeout",
		Status:          "paused",
		Branch:          "fix/payment-timeout",
		BaseBranch:      "main",
		Dirty:           true,
		ChangedFiles:    []string{"src/payment/retry-policy.ts"},
		CurrentNote:     "Fix duplicate retries",
		Services:        []string{"api", "frontend"},
		LastCheckCmd:    "pnpm test payment-retry",
		LastCheckResult: "failed",
		LastCheckExit:   1,
	}

	md := GenerateHandoff(data)

	if !strings.Contains(md, "Handoff: payment-timeout") {
		t.Error("missing title")
	}
	if !strings.Contains(md, "fix/payment-timeout") {
		t.Error("missing branch")
	}
	if !strings.Contains(md, "Fix duplicate retries") {
		t.Error("missing note")
	}
	if !strings.Contains(md, "No environment variable values are included") {
		t.Error("missing security notice")
	}
}

func TestRedact(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Bearer abc123def456", "Bearer [REDACTED]"},
		{"password=supersecret", "password=[REDACTED]"},
		{"token=abc123", "token=[REDACTED]"},
		{"plain text no secrets", "plain text no secrets"},
	}

	for _, tt := range tests {
		got := Redact(tt.input)
		if got != tt.expected {
			t.Errorf("Redact(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
