package library

import (
	"testing"

	"github.com/matryer/is"
)

func TestExtractSuggestedID(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "file path",
			source: "/path/to/my-schema.graphqls",
			want:   "my-schema",
		},
		{
			name:   "file path with spaces",
			source: "/path/to/My Schema.graphqls",
			want:   "my-schema",
		},
		{
			name:   "URL with api prefix",
			source: "https://api.github.com/graphql",
			want:   "github",
		},
		{
			name:   "URL without api prefix",
			source: "https://shopify.com/graphql",
			want:   "shopify",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got := ExtractSuggestedID(tt.source)
			is.Equal(got, tt.want)
		})
	}
}
