package assert

import (
	"strings"
	"testing"
)

type Assert struct {
	t *testing.T
}

func New(t *testing.T) Assert {
	return Assert{t: t}
}

// StringContains checks if text contains substring and prints both if it doesn't
func (a *Assert) StringContains(text, substring string) {
	a.t.Helper()
	if !strings.Contains(text, substring) {
		a.t.Fatalf("Expected text to contain substring\n\nActual text:\n%s\n\nExpected substring:\n%s\n", text, substring)
	}
}
