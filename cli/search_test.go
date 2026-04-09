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

func TestUnknownFieldWarnings(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		wantMsgs []string
	}{
		{
			name:     "no field specifiers",
			query:    "user",
			wantMsgs: nil,
		},
		{
			name:     "field with space after colon is not a field specifier",
			query:    "kind: Query",
			wantMsgs: nil,
		},
		{
			name:     "valid field",
			query:    "kind:Query",
			wantMsgs: nil,
		},
		{
			name:     "all valid fields",
			query:    "+kind:Query +name:user +usage:User",
			wantMsgs: nil,
		},
		{
			name:     "unknown field with suggestion",
			query:    "naem:user",
			wantMsgs: []string{`⚠️ unknown search field "naem" — did you mean: name?`},
		},
		{
			name:     "unknown field no suggestion",
			query:    "foo:bar",
			wantMsgs: []string{`⚠️ unknown search field "foo"`},
		},
		{
			name:     "typo in kind",
			query:    "knd:Query",
			wantMsgs: []string{`⚠️ unknown search field "knd" — did you mean: kind?`},
		},
		{
			name:     "typo in description",
			query:    "descripton:hello",
			wantMsgs: []string{`⚠️ unknown search field "descripton" — did you mean: description?`},
		},
		{
			name:     "multiple unknown fields",
			query:    "knd:Query naem:user",
			wantMsgs: []string{`⚠️ unknown search field "knd" — did you mean: kind?`, `⚠️ unknown search field "naem" — did you mean: name?`},
		},
		{
			name:     "duplicate unknown field only warned once",
			query:    "foo:a foo:b",
			wantMsgs: []string{`⚠️ unknown search field "foo"`},
		},
		{
			name:     "field name casing must match exactly",
			query:    "Kind:Query",
			wantMsgs: []string{`⚠️ unknown search field "Kind" — did you mean: kind?`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got := unknownFieldWarnings(tt.query)
			is.Equal(got, tt.wantMsgs)
		})
	}
}

func TestApplyKindFilter(t *testing.T) {
	tests := []struct {
		name       string
		kindFilter string
		query      string
		wantQuery  string
		wantErr    string
	}{
		{
			name:       "no query",
			kindFilter: "Query",
			query:      "",
			wantQuery:  "+kind:Query",
		},
		{
			name:       "with query",
			kindFilter: "Object",
			query:      "user",
			wantQuery:  "+kind:Object user",
		},
		{
			name:       "case insensitive",
			kindFilter: "object",
			query:      "",
			wantQuery:  "+kind:Object",
		},
		{
			name:       "mixed case",
			kindFilter: "MUTATION",
			query:      "",
			wantQuery:  "+kind:Mutation",
		},
		{
			name:       "ObjectField",
			kindFilter: "objectfield",
			query:      "email",
			wantQuery:  "+kind:ObjectField email",
		},
		{
			name:       "invalid kind",
			kindFilter: "BadType",
			query:      "",
			wantErr:    `invalid --kind value "BadType"`,
		},
		{
			name:       "conflict with kind: in query",
			kindFilter: "Query",
			query:      "+kind:Object user",
			wantErr:    "cannot use --kind flag when query already contains a kind: filter",
		},
		{
			name:       "conflict with bare kind: in query",
			kindFilter: "Query",
			query:      "kind:Object",
			wantErr:    "cannot use --kind flag when query already contains a kind: filter",
		},
		{
			name:       "conflict with -kind: exclusion in query",
			kindFilter: "Query",
			query:      "-kind:*Field",
			wantErr:    "cannot use --kind flag when query already contains a kind: filter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got, err := applyKindFilter(tt.kindFilter, tt.query)
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
					Kind:        "Object",
					Name:        "User",
					Path:        "User",
					Description: "A user",
					Score:       1.5,
				},
			},
			want: []map[string]any{
				{
					"kind":        "Object",
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
					Kind:        "Query",
					Name:        "getUser",
					Path:        "Query.getUser",
					Description: "Get a user by ID",
					Score:       1.5,
					Signature:   "getUser(id: ID!): User",
				},
			},
			want: []map[string]any{
				{
					"kind":        "Query",
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
					Kind:        "Object",
					Name:        "User",
					Path:        "User",
					Description: "A user",
					Score:       2.0,
				},
				{
					Kind:        "Query",
					Name:        "userId",
					Path:        "Query.userId",
					Description: "Get user by ID",
					Score:       1.5,
					Signature:   "userId(id: ID!): User",
				},
			},
			want: []map[string]any{
				{
					"kind":        "Object",
					"name":        "User",
					"path":        "User",
					"description": "A user",
					"score":       2.0,
					"signature":   "",
				},
				{
					"kind":        "Query",
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
				is.Equal(got[i]["kind"], wantItem["kind"])
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
			Kind:        "Object",
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
