package capsule

import (
	"fmt"
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$`)

func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("capsule name must not be empty")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("capsule name must not be %q", name)
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("capsule name must not contain path separators")
	}
	if len(name) > 64 {
		return fmt.Errorf("capsule name must be at most 64 characters")
	}
	if !slugPattern.MatchString(name) {
		return fmt.Errorf("capsule name must start with a letter or digit and contain only letters, digits, dashes, or underscores")
	}
	return nil
}

func Slugify(name string) string {
	return name
}
