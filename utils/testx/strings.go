package testx

import (
	"regexp"
	"strings"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

// StripANSI removes ANSI escape codes from a string.
func StripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// NormalizeView strips ANSI escape codes, leading/trailing whitespace, and empty lines from multi-line strings.
// This is useful for testing bubble-tea views which may have ANSI color codes and empty-lines for spacing.
func NormalizeView(s string) string {
	s = StripANSI(s)
	lines := strings.Split(s, "\n")
	var result []string

	// Trim whitespace from each line and skip empty lines
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return strings.Join(result, "\n")
}
