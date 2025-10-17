package testx

import "strings"

// NormalizeView strips leading/trailing whitespace and empty lines from multi-line strings
// This is useful for testing bubble-tea views which may have empty-lines for spacing.
func NormalizeView(s string) string {
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
