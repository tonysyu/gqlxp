package text

import (
	"testing"

	"github.com/matryer/is"
)

func TestSuggestSimilarWords(t *testing.T) {
	allowed := []string{"kind", "name", "description", "path", "usage", "implements"}

	tests := []struct {
		word string
		want []string
	}{
		{"kind", []string{"kind"}},              // exact match
		{"knd", []string{"kind"}},               // 1 deletion
		{"naem", []string{"name"}},              // transposition
		{"pth", []string{"path"}},               // 1 deletion
		{"usag", []string{"usage"}},             // 1 deletion
		{"descripton", []string{"description"}}, // 1 deletion
		{"implments", []string{"implements"}},   // 1 deletion
		{"foo", nil},                            // no close match
		{"xyz", nil},                            // no close match
		{"KIND", []string{"kind"}},              // case-insensitive
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			is := is.New(t)
			is.Equal(SuggestSimilarWords(tt.word, allowed), tt.want)
		})
	}
}
