package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/search"
)

func TestApplyTypeFilter(t *testing.T) {
	tests := []struct {
		name       string
		typeFilter string
		query      string
		wantQuery  string
		wantErr    string
	}{
		{
			name:       "no query",
			typeFilter: "Query",
			query:      "",
			wantQuery:  "+type:Query",
		},
		{
			name:       "with query",
			typeFilter: "Object",
			query:      "user",
			wantQuery:  "+type:Object user",
		},
		{
			name:       "case insensitive",
			typeFilter: "object",
			query:      "",
			wantQuery:  "+type:Object",
		},
		{
			name:       "mixed case",
			typeFilter: "MUTATION",
			query:      "",
			wantQuery:  "+type:Mutation",
		},
		{
			name:       "ObjectField",
			typeFilter: "objectfield",
			query:      "email",
			wantQuery:  "+type:ObjectField email",
		},
		{
			name:       "invalid type",
			typeFilter: "BadType",
			query:      "",
			wantErr:    `invalid --type value "BadType"`,
		},
		{
			name:       "conflict with type: in query",
			typeFilter: "Query",
			query:      "+type:Object user",
			wantErr:    "cannot use --type flag when query already contains a type: filter",
		},
		{
			name:       "conflict with bare type: in query",
			typeFilter: "Query",
			query:      "type:Object",
			wantErr:    "cannot use --type flag when query already contains a type: filter",
		},
		{
			name:       "conflict with -type: exclusion in query",
			typeFilter: "Query",
			query:      "-type:*Field",
			wantErr:    "cannot use --type flag when query already contains a type: filter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got, err := applyTypeFilter(tt.typeFilter, tt.query)
			if tt.wantErr != "" {
				is.True(err != nil)
				is.True(strings.Contains(err.Error(), tt.wantErr))
			} else {
				is.NoErr(err)
				is.Equal(got, tt.wantQuery)
			}
		})
	}
}

func TestPrintSearchResultsJSON(t *testing.T) {
	tests := []struct {
		name    string
		results []search.SearchResult
		want    []map[string]any
	}{
		{
			name:    "empty results",
			results: []search.SearchResult{},
			want:    []map[string]any{},
		},
		{
			name: "single result without signature",
			results: []search.SearchResult{
				{
					Type:        "Object",
					Name:        "User",
					Path:        "User",
					Description: "A user",
					Score:       1.5,
				},
			},
			want: []map[string]any{
				{
					"type":        "Object",
					"name":        "User",
					"path":        "User",
					"description": "A user",
					"score":       1.5,
					"signature":   "",
				},
			},
		},
		{
			name: "result with signature",
			results: []search.SearchResult{
				{
					Type:        "Query",
					Name:        "getUser",
					Path:        "Query.getUser",
					Description: "Get a user by ID",
					Score:       1.5,
					Signature:   "getUser(id: ID!): User",
				},
			},
			want: []map[string]any{
				{
					"type":        "Query",
					"name":        "getUser",
					"path":        "Query.getUser",
					"description": "Get a user by ID",
					"score":       1.5,
					"signature":   "getUser(id: ID!): User",
				},
			},
		},
		{
			name: "multiple results",
			results: []search.SearchResult{
				{
					Type:        "Object",
					Name:        "User",
					Path:        "User",
					Description: "A user",
					Score:       2.0,
				},
				{
					Type:        "Query",
					Name:        "userId",
					Path:        "Query.userId",
					Description: "Get user by ID",
					Score:       1.5,
					Signature:   "userId(id: ID!): User",
				},
			},
			want: []map[string]any{
				{
					"type":        "Object",
					"name":        "User",
					"path":        "User",
					"description": "A user",
					"score":       2.0,
					"signature":   "",
				},
				{
					"type":        "Query",
					"name":        "userId",
					"path":        "Query.userId",
					"description": "Get user by ID",
					"score":       1.5,
					"signature":   "userId(id: ID!): User",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := printSearchResultsJSON(tt.results)
			is.NoErr(err)

			// Restore stdout and read captured output
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			is.NoErr(err)
			output := buf.String()

			// Parse the JSON output
			var got []map[string]any
			err = json.Unmarshal([]byte(output), &got)
			is.NoErr(err)

			// Verify the output matches expected
			is.Equal(len(got), len(tt.want))
			for i, wantItem := range tt.want {
				is.Equal(got[i]["type"], wantItem["type"])
				is.Equal(got[i]["name"], wantItem["name"])
				is.Equal(got[i]["path"], wantItem["path"])
				is.Equal(got[i]["description"], wantItem["description"])
				is.Equal(got[i]["score"], wantItem["score"])
				is.Equal(got[i]["signature"], wantItem["signature"])
			}
		})
	}
}

func TestPrintSearchResultsJSON_PrettyPrinted(t *testing.T) {
	is := is.New(t)

	results := []search.SearchResult{
		{
			Type:        "Object",
			Name:        "User",
			Path:        "User",
			Description: "A user",
			Score:       1.5,
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := printSearchResultsJSON(results)
	is.NoErr(err)

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	is.NoErr(err)
	output := buf.String()

	// Verify output is pretty-printed (contains 2-space indentation)
	is.True(len(output) > 0)
	is.True(bytes.Contains([]byte(output), []byte("  "))) // Contains indentation
	is.True(bytes.Contains([]byte(output), []byte("\n"))) // Contains newlines
}
