package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tonysyu/gqlxp/cli/prompt"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/library"
)

// Prompter abstracts interactive terminal prompts for schema resolution.
type Prompter interface {
	YesNo(prompt string) (bool, error)
	SchemaID(suggested string) (string, error)
	String(prompt, defaultValue string) (string, error)
}

// terminalPrompter is the real implementation that delegates to the prompt package.
type terminalPrompter struct{}

func (terminalPrompter) YesNo(p string) (bool, error) {
	return prompt.YesNo(p)
}

func (terminalPrompter) SchemaID(suggested string) (string, error) {
	return prompt.SchemaID(suggested)
}

func (terminalPrompter) String(p, defaultValue string) (string, error) {
	return prompt.String(p, defaultValue)
}

// LoadedSchema is the ready-to-use result of schema resolution.
type LoadedSchema struct {
	ID        string
	Content   []byte
	GQLSchema gql.GraphQLSchema
}

// SchemaLoader resolves a CLI schema argument into a LoadedSchema.
type SchemaLoader struct {
	lib      library.Library
	prompter Prompter
}

// NewSchemaLoader creates a SchemaLoader with injected dependencies.
func NewSchemaLoader(lib library.Library, p Prompter) *SchemaLoader {
	return &SchemaLoader{lib: lib, prompter: p}
}

// NewDefaultSchemaLoader creates a SchemaLoader backed by the real filesystem library
// and real terminal prompter. This is the one-liner for CLI commands.
func NewDefaultSchemaLoader() *SchemaLoader {
	return NewSchemaLoader(library.NewLibrary(), terminalPrompter{})
}

// LoadSchema is a convenience wrapper around NewDefaultSchemaLoader().Load(arg).
func LoadSchema(arg string) (LoadedSchema, error) {
	return NewDefaultSchemaLoader().Load(arg)
}

// Load resolves a schema argument and returns a LoadedSchema with a parsed schema.
// arg can be:
//   - Empty string: use default schema from config
//   - A schema ID that exists in the library
//   - A file path (will be added to library if needed)
func (l *SchemaLoader) Load(arg string) (LoadedSchema, error) {
	var schemaID string
	var content []byte

	if arg == "" {
		defaultSchemaID, err := l.lib.GetDefaultSchema()
		if err != nil {
			return LoadedSchema{}, fmt.Errorf("error getting default schema: %w", err)
		}
		if defaultSchemaID == "" {
			return LoadedSchema{}, fmt.Errorf("no schema specified and no default schema set. Use 'gqlxp library default' to set one")
		}
		schemaID = defaultSchemaID
	} else {
		// First check if it's an existing schema ID
		if _, err := l.lib.Get(arg); err == nil {
			schemaID = arg
		} else {
			// Not a schema ID - try as file path
			resolvedID, resolvedContent, err := l.resolveFilePath(arg)
			if err != nil {
				return LoadedSchema{}, fmt.Errorf("invalid schema argument '%s': %w", arg, err)
			}
			schemaID = resolvedID
			content = resolvedContent
		}
	}

	// Load content from library if not already resolved from file path
	if content == nil {
		schema, err := l.lib.Get(schemaID)
		if err != nil {
			return LoadedSchema{}, fmt.Errorf("failed to load schema '%s': %w", schemaID, err)
		}
		content = schema.Content
	}

	parsedSchema, err := gql.ParseSchema(content)
	if err != nil {
		return LoadedSchema{}, fmt.Errorf("error parsing schema: %w", err)
	}

	return LoadedSchema{ID: schemaID, Content: content, GQLSchema: parsedSchema}, nil
}

// resolveFilePath handles the case where arg is a file path.
func (l *SchemaLoader) resolveFilePath(filePath string) (schemaID string, content []byte, err error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	content, err = os.ReadFile(absPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	fileHash := library.CalculateFileHash(content)
	existingSchema, err := l.lib.FindByPath(absPath)

	// No match - register new schema
	if err != nil {
		id, err := l.registerSchema(absPath, content)
		return id, content, err
	}

	// Hash matches - use existing schema
	if existingSchema.Metadata.FileHash == fileHash {
		return existingSchema.ID, existingSchema.Content, nil
	}

	// Hash mismatch - handle update workflow
	return l.handleSchemaUpdate(existingSchema, content)
}

func (l *SchemaLoader) handleSchemaUpdate(existingSchema *library.Schema, newContent []byte) (string, []byte, error) {
	fmt.Printf("Schema file has changed since last import.\n")
	update, err := l.prompter.YesNo("Update library")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user input: %w", err)
	}

	if !update {
		fmt.Printf("Using existing library version\n")
		return existingSchema.ID, existingSchema.Content, nil
	}

	if err := l.lib.UpdateContent(existingSchema.ID, newContent); err != nil {
		return "", nil, fmt.Errorf("failed to update library: %w", err)
	}

	fmt.Printf("Library schema '%s' updated\n", existingSchema.ID)
	return existingSchema.ID, newContent, nil
}

func (l *SchemaLoader) registerSchema(filePath string, content []byte) (string, error) {
	basename := filepath.Base(filePath)
	ext := filepath.Ext(basename)
	suggested := library.SanitizeSchemaID(basename[:len(basename)-len(ext)])

	schemaID, err := l.prompter.SchemaID(suggested)
	if err != nil {
		return "", fmt.Errorf("failed to get schema ID: %w", err)
	}

	displayName, err := l.prompter.String("Enter display name", schemaID)
	if err != nil {
		return "", fmt.Errorf("failed to get display name: %w", err)
	}

	if err := l.lib.Add(schemaID, displayName, filePath); err != nil {
		return "", fmt.Errorf("failed to add schema to library: %w", err)
	}

	fmt.Printf("Schema '%s' added to library\n", schemaID)
	return schemaID, nil
}
