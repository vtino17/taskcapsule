package report

import (
	"regexp"
	"strings"
)

var redactPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(Bearer\s+)[A-Za-z0-9._\-+/=]{10,}`),
	regexp.MustCompile(`(?i)(Authorization:\s*)[A-Za-z0-9._\-+/=]{10,}`),
	regexp.MustCompile(`(?i)(password\s*=\s*)[^\s&"]{4,}`),
	regexp.MustCompile(`(?i)(token\s*=\s*)[^\s&"]{4,}`),
	regexp.MustCompile(`(?i)(api_key\s*=\s*)[^\s&"]{4,}`),
	regexp.MustCompile(`(?i)(secret\s*=\s*)[^\s&"]{4,}`),
	regexp.MustCompile(`(?i)(-----BEGIN\s+(?:RSA\s+)?PRIVATE\s+KEY-----)`),
}

var redactPlaceholder = []byte("$1[REDACTED]")

func Redact(content string) string {
	result := []byte(content)
	for _, pattern := range redactPatterns {
		result = pattern.ReplaceAll(result, redactPlaceholder)
	}
	return string(result)
}

func ContainsSecret(content string) bool {
	for _, pattern := range redactPatterns {
		if pattern.MatchString(content) {
			return true
		}
	}

	lower := strings.ToLower(content)
	sensitiveTerms := []string{"api_key", "apikey", "secret", "password", "token", "jwt", "bearer "}
	for _, term := range sensitiveTerms {
		if strings.Contains(lower, term) {
			return true
		}
	}

	return false
}
