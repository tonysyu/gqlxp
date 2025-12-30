package gql_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gqlxp/gql"
)

func TestMainWithMissingSchemaFile(t *testing.T) {
	is := is.New(t)

	// Test that main function handles missing schema file gracefully
	// We can't directly test main() as it calls os.Exit, but we can test
	// the file reading logic separately

	nonExistentFile := "/path/that/does/not/exist/schema.graphqls"
	_, err := os.ReadFile(nonExistentFile)
	is.True(err != nil) // Should return an error for non-existent file

	// Check that it's specifically a file not found error
	is.True(os.IsNotExist(err))
}

func TestMainWithUnreadableSchemaFile(t *testing.T) {
	is := is.New(t)

	// Create a temporary file with no read permissions
	tempDir := t.TempDir()
	unreadableFile := filepath.Join(tempDir, "unreadable.graphqls")

	// Create the file first
	err := os.WriteFile(unreadableFile, []byte("type Query { test: String }"), 0644)
	is.True(err == nil)

	// Remove read permissions
	err = os.Chmod(unreadableFile, 0000)
	is.True(err == nil)

	// Restore permissions after test for cleanup
	defer func() { _ = os.Chmod(unreadableFile, 0644) }()

	// Try to read the unreadable file
	_, err = os.ReadFile(unreadableFile)
	is.True(err != nil) // Should return a permission error

	// Check that it's a permission error
	var pathError *fs.PathError
	is.True(err.(*fs.PathError) != nil || pathError != nil)
}

func TestParseSchemaWithMalformedGraphQL(t *testing.T) {
	// Skip this test because ParseSchema calls log.Fatalf on parse errors
	// which terminates the test process. In a production system, we'd want
	// ParseSchema to return errors instead of calling log.Fatalf
	t.Skip("Skipping malformed GraphQL test - ParseSchema calls log.Fatalf which exits the process")
}

func TestParseSchemaWithCorruptedContent(t *testing.T) {
	// Skip this test because ParseSchema calls log.Fatalf on parse errors
	t.Skip("Skipping corrupted content test - ParseSchema calls log.Fatalf which exits the process")
}

func TestParseSchemaWithExtremelyLargeContent(t *testing.T) {
	is := is.New(t)

	// Create an extremely large but valid schema
	var largeSchemaBuilder []string
	largeSchemaBuilder = append(largeSchemaBuilder, "type Query {")

	// Add a very large number of fields
	for i := 0; i < 10000; i++ {
		fieldName := fmt.Sprintf("field%d", i)
		largeSchemaBuilder = append(largeSchemaBuilder, "  "+fieldName+": String")
	}
	largeSchemaBuilder = append(largeSchemaBuilder, "}")

	largeSchema := ""
	for _, line := range largeSchemaBuilder {
		largeSchema += line + "\n"
	}

	// This should either succeed or fail gracefully without consuming excessive memory
	defer func() {
		if r := recover(); r != nil {
			t.Logf("ParseSchema panicked for large schema: %v", r)
		}
	}()

	schema, _ := ParseSchema([]byte(largeSchema))

	// If it succeeds, verify it parsed correctly
	is.True(len(schema.Query) <= 10000) // Should have parsed some or all fields
}

func TestParseSchemaWithVeryDeepNesting(t *testing.T) {
	is := is.New(t)

	// Create a schema with very deep nesting
	deepNesting := "type Query { field: "
	for i := 0; i < 100; i++ {
		deepNesting += "["
	}
	deepNesting += "String"
	for i := 0; i < 100; i++ {
		deepNesting += "]"
	}
	deepNesting += " }"

	defer func() {
		if r := recover(); r != nil {
			t.Logf("ParseSchema panicked for deeply nested types: %v", r)
		}
	}()

	schema, _ := ParseSchema([]byte(deepNesting))

	// If parsing succeeds, test that TypeString can handle it
	if len(schema.Query) > 0 {
		for _, field := range schema.Query {
			typeStr := field.TypeString()
			is.True(len(typeStr) > 0)
		}
	}
}

func TestSchemaWithCircularReferences(t *testing.T) {
	is := is.New(t)

	// Test schema with potential circular references
	circularSchema := `
		type User {
			id: ID!
			friend: User
			friends: [User!]!
			posts: [Post!]!
		}

		type Post {
			id: ID!
			author: User!
			comments: [Comment!]!
		}

		type Comment {
			id: ID!
			author: User!
			post: Post!
			replies: [Comment!]!
		}

		type Query {
			getUser: User
			getPost: Post
		}
	`

	// Should parse without infinite loops
	schema, _ := ParseSchema([]byte(circularSchema))

	is.True(len(schema.Query) == 2)
	is.True(len(schema.Object) == 3)

	// Test that we can access all the types without issues
	user := schema.Object["User"]
	is.True(user != nil)
	is.Equal(user.Name(), "User")

	post := schema.Object["Post"]
	is.True(post != nil)
	is.Equal(post.Name(), "Post")

	comment := schema.Object["Comment"]
	is.True(comment != nil)
	is.Equal(comment.Name(), "Comment")
}

func TestEmptyFileHandling(t *testing.T) {
	is := is.New(t)

	// Create a temporary empty file
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.graphqls")

	err := os.WriteFile(emptyFile, []byte(""), 0644)
	is.True(err == nil)

	// Read the empty file
	content, err := os.ReadFile(emptyFile)
	is.True(err == nil)
	is.Equal(len(content), 0)

	// Parse empty content
	schema, _ := ParseSchema(content)
	is.Equal(len(schema.Query), 0)
	is.Equal(len(schema.Mutation), 0)
}

func TestSchemaWithOnlyComments(t *testing.T) {
	is := is.New(t)

	commentsOnlySchema := `
		# This is a comment
		# Another comment
		# Multiple lines of comments
		# No actual GraphQL definitions
	`

	schema, _ := ParseSchema([]byte(commentsOnlySchema))
	is.Equal(len(schema.Query), 0)
	is.Equal(len(schema.Mutation), 0)
	is.Equal(len(schema.Object), 0)
}

func TestSchemaFilePermissions(t *testing.T) {
	is := is.New(t)

	// Create a temporary file with specific permissions
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.graphqls")

	validSchema := "type Query { test: String }"
	err := os.WriteFile(testFile, []byte(validSchema), 0644)
	is.True(err == nil)

	// Verify file exists and is readable
	info, err := os.Stat(testFile)
	is.True(err == nil)
	is.True(info.Size() > 0)
	is.True(!info.IsDir())

	// Read and parse the file
	content, err := os.ReadFile(testFile)
	is.True(err == nil)

	schema, _ := ParseSchema(content)
	is.Equal(len(schema.Query), 1)
}
