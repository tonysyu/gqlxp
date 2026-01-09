package cli

import (
	"testing"

	"github.com/matryer/is"
)

func TestParseSelectionTarget(t *testing.T) {
	tests := []struct {
		name          string
		target        string
		wantTypeName  string
		wantFieldName string
	}{
		{
			name:          "empty string",
			target:        "",
			wantTypeName:  "",
			wantFieldName: "",
		},
		{
			name:          "type only",
			target:        "User",
			wantTypeName:  "User",
			wantFieldName: "",
		},
		{
			name:          "type and field",
			target:        "Query.user",
			wantTypeName:  "Query",
			wantFieldName: "user",
		},
		{
			name:          "multiple dots - only first split",
			target:        "A.B.C",
			wantTypeName:  "A",
			wantFieldName: "B.C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			gotTypeName, gotFieldName := parseSelectionTarget(tt.target)
			is.Equal(gotTypeName, tt.wantTypeName)   // parseSelectionTarget() typeName
			is.Equal(gotFieldName, tt.wantFieldName) // parseSelectionTarget() fieldName
		})
	}
}
