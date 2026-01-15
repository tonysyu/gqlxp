package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/search"
)

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
			name: "single result",
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
					Type:        "Field",
					Name:        "userId",
					Path:        "Query.userId",
					Description: "Get user by ID",
					Score:       1.5,
				},
			},
			want: []map[string]any{
				{
					"type":        "Object",
					"name":        "User",
					"path":        "User",
					"description": "A user",
					"score":       2.0,
				},
				{
					"type":        "Field",
					"name":        "userId",
					"path":        "Query.userId",
					"description": "Get user by ID",
					"score":       1.5,
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
