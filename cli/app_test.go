package cli

import "testing"

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
			gotTypeName, gotFieldName := parseSelectionTarget(tt.target)
			if gotTypeName != tt.wantTypeName {
				t.Errorf("parseSelectionTarget() typeName = %v, want %v", gotTypeName, tt.wantTypeName)
			}
			if gotFieldName != tt.wantFieldName {
				t.Errorf("parseSelectionTarget() fieldName = %v, want %v", gotFieldName, tt.wantFieldName)
			}
		})
	}
}
