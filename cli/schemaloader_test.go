package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/library"
)

// fakeLib is a minimal in-memory Library for testing SchemaLoader.
type fakeLib struct {
	schemas       map[string]*library.Schema
	defaultSchema string
	// Track calls for assertion
	addCalled     bool
	updateCalled  bool
	updateID      string
	updateContent []byte
}

func newFakeLib() *fakeLib {
	return &fakeLib{schemas: make(map[string]*library.Schema)}
}

func (f *fakeLib) withSchema(id string, content []byte, hash string, sourcePath string) *fakeLib {
	f.schemas[id] = &library.Schema{
		ID:      id,
		Content: content,
		Metadata: library.SchemaMetadata{
			FileHash:    hash,
			SourceFile:  sourcePath,
			DisplayName: id,
			URLPatterns: make(map[string]string),
		},
	}
	return f
}

func (f *fakeLib) withDefault(id string) *fakeLib {
	f.defaultSchema = id
	return f
}

func (f *fakeLib) Get(id string) (*library.Schema, error) {
	s, ok := f.schemas[id]
	if !ok {
		return nil, errors.New("schema not found: " + id)
	}
	return s, nil
}

func (f *fakeLib) GetDefaultSchema() (string, error) {
	return f.defaultSchema, nil
}

func (f *fakeLib) FindByPath(absolutePath string) (*library.Schema, error) {
	for _, s := range f.schemas {
		if s.Metadata.SourceFile == absolutePath {
			return s, nil
		}
	}
	return nil, errors.New("no schema for path: " + absolutePath)
}

func (f *fakeLib) UpdateContent(id string, content []byte) error {
	f.updateCalled = true
	f.updateID = id
	f.updateContent = content
	if s, ok := f.schemas[id]; ok {
		s.Content = content
		s.Metadata.FileHash = library.CalculateFileHash(content)
	}
	return nil
}

func (f *fakeLib) Add(id, displayName, sourcePath string) error {
	f.addCalled = true
	content, _ := os.ReadFile(sourcePath)
	f.schemas[id] = &library.Schema{
		ID:      id,
		Content: content,
		Metadata: library.SchemaMetadata{
			DisplayName: displayName,
			SourceFile:  sourcePath,
			FileHash:    library.CalculateFileHash(content),
			URLPatterns: make(map[string]string),
		},
	}
	return nil
}

// Unused interface methods.
func (f *fakeLib) AddFromContent(id, displayName string, content []byte, sourceInfo string) error {
	return nil
}
func (f *fakeLib) List() ([]library.SchemaInfo, error)                             { return nil, nil }
func (f *fakeLib) Remove(id string) error                                          { return nil }
func (f *fakeLib) UpdateMetadata(id string, metadata library.SchemaMetadata) error { return nil }
func (f *fakeLib) SetURLPattern(id, typePattern, urlPattern string) error          { return nil }
func (f *fakeLib) SetDefaultSchema(id string) error                                { return nil }
func (f *fakeLib) EnsureIndex(schemaID string, schema *gql.GraphQLSchema) error    { return nil }
func (f *fakeLib) Reindex(schemaID string) error                                   { return nil }

// fakePrompter records calls and returns pre-configured answers.
type fakePrompter struct {
	yesNoResult    bool
	schemaIDResult string
	yesNoCalled    bool
	schemaIDCalled bool
}

func (f *fakePrompter) YesNo(_ string) (bool, error) {
	f.yesNoCalled = true
	return f.yesNoResult, nil
}

func (f *fakePrompter) SchemaID(_ string) (string, error) {
	f.schemaIDCalled = true
	return f.schemaIDResult, nil
}

func (f *fakePrompter) String(_, defaultValue string) (string, error) {
	return defaultValue, nil
}

const validSchema = `type Query { hello: String }`

func writeSchemaFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.graphqls")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSchemaLoader_EmptyArg_WithDefault(t *testing.T) {
	is := is.New(t)

	lib := newFakeLib().
		withSchema("my-schema", []byte(validSchema), "hash", "").
		withDefault("my-schema")
	loader := NewSchemaLoader(lib, &fakePrompter{})

	schema, err := loader.Load("")

	is.NoErr(err)
	is.Equal(schema.ID, "my-schema")
}

func TestSchemaLoader_EmptyArg_NoDefault(t *testing.T) {
	is := is.New(t)

	loader := NewSchemaLoader(newFakeLib(), &fakePrompter{})

	_, err := loader.Load("")

	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "no default schema"))
}

func TestSchemaLoader_KnownSchemaID(t *testing.T) {
	is := is.New(t)

	lib := newFakeLib().withSchema("github", []byte(validSchema), "hash", "")
	loader := NewSchemaLoader(lib, &fakePrompter{})

	schema, err := loader.Load("github")

	is.NoErr(err)
	is.Equal(schema.ID, "github")
}

func TestSchemaLoader_FilePath_HashMatches(t *testing.T) {
	is := is.New(t)

	path := writeSchemaFile(t, validSchema)
	absPath, _ := filepath.Abs(path)
	hash := library.CalculateFileHash([]byte(validSchema))
	prompter := &fakePrompter{}

	lib := newFakeLib().withSchema("existing", []byte(validSchema), hash, absPath)
	loader := NewSchemaLoader(lib, prompter)

	schema, err := loader.Load(path)

	is.NoErr(err)
	is.Equal(schema.ID, "existing")
	is.True(!prompter.yesNoCalled) // no prompt shown when hash matches
}

func TestSchemaLoader_FilePath_HashMismatch_Update(t *testing.T) {
	is := is.New(t)

	newContent := validSchema + "\ntype User { id: ID! }"
	path := writeSchemaFile(t, newContent)
	absPath, _ := filepath.Abs(path)
	oldHash := library.CalculateFileHash([]byte(validSchema))
	prompter := &fakePrompter{yesNoResult: true}

	lib := newFakeLib().withSchema("existing", []byte(validSchema), oldHash, absPath)
	loader := NewSchemaLoader(lib, prompter)

	schema, err := loader.Load(path)

	is.NoErr(err)
	is.Equal(schema.ID, "existing")
	is.True(prompter.yesNoCalled) // prompted about update
	is.True(lib.updateCalled)     // library was updated
}

func TestSchemaLoader_FilePath_HashMismatch_KeepExisting(t *testing.T) {
	is := is.New(t)

	newContent := validSchema + "\ntype User { id: ID! }"
	path := writeSchemaFile(t, newContent)
	absPath, _ := filepath.Abs(path)
	oldHash := library.CalculateFileHash([]byte(validSchema))
	prompter := &fakePrompter{yesNoResult: false}

	lib := newFakeLib().withSchema("existing", []byte(validSchema), oldHash, absPath)
	loader := NewSchemaLoader(lib, prompter)

	schema, err := loader.Load(path)

	is.NoErr(err)
	is.Equal(schema.ID, "existing")
	is.True(prompter.yesNoCalled) // prompted about update
	is.True(!lib.updateCalled)    // library was NOT updated
}

func TestSchemaLoader_FilePath_NewFile_Registers(t *testing.T) {
	is := is.New(t)

	path := writeSchemaFile(t, validSchema)
	prompter := &fakePrompter{schemaIDResult: "new-schema"}

	lib := newFakeLib()
	loader := NewSchemaLoader(lib, prompter)

	schema, err := loader.Load(path)

	is.NoErr(err)
	is.Equal(schema.ID, "new-schema")
	is.True(prompter.schemaIDCalled) // prompted for schema ID
	is.True(lib.addCalled)           // schema was added to library
}

func TestSchemaLoader_ParsedSchemaIsValid(t *testing.T) {
	is := is.New(t)

	lib := newFakeLib().withSchema("test", []byte(validSchema), "hash", "")
	loader := NewSchemaLoader(lib, &fakePrompter{})

	schema, err := loader.Load("test")

	is.NoErr(err)
	// Verify the schema was actually parsed — Query type should exist
	is.True(schema.GQLSchema.Query != nil)
}
