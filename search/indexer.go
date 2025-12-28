package search

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/tonysyu/gqlxp/gql"
)

// document represents a searchable item in the schema
type document struct {
	Type        string `json:"type"`        // Object, Field, Enum, etc.
	Name        string `json:"name"`        // Type or field name
	Description string `json:"description"` // Description text
	Path        string `json:"path"`        // Full path (e.g., "Query.user.name")
	SchemaID    string `json:"schemaID"`    // Schema identifier
}

// BleveIndexer implements Indexer using Bleve
type BleveIndexer struct {
	baseDir string // Base directory for storing indexes
}

// NewIndexer creates a new BleveIndexer with the given base directory
func NewIndexer(baseDir string) *BleveIndexer {
	return &BleveIndexer{baseDir: baseDir}
}

// Index creates or updates the index for a schema
func (b *BleveIndexer) Index(schemaID string, schema *gql.GraphQLSchema) error {
	indexPath := b.getIndexPath(schemaID)

	// Remove existing index if it exists
	if err := os.RemoveAll(indexPath); err != nil {
		return fmt.Errorf("failed to remove existing index: %w", err)
	}

	// Create new index
	index, err := bleve.New(indexPath, buildIndexMapping())
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer index.Close()

	// Extract and index documents
	docs := extractDocuments(schemaID, schema)
	batch := index.NewBatch()

	for i, doc := range docs {
		id := fmt.Sprintf("%s-%d", schemaID, i)
		if err := batch.Index(id, doc); err != nil {
			return fmt.Errorf("failed to add document to batch: %w", err)
		}
	}

	if err := index.Batch(batch); err != nil {
		return fmt.Errorf("failed to index batch: %w", err)
	}

	return nil
}

// Remove deletes the index for a schema
func (b *BleveIndexer) Remove(schemaID string) error {
	indexPath := b.getIndexPath(schemaID)
	if err := os.RemoveAll(indexPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove index: %w", err)
	}
	return nil
}

// Exists checks if an index exists for a schema
func (b *BleveIndexer) Exists(schemaID string) bool {
	indexPath := b.getIndexPath(schemaID)
	_, err := os.Stat(indexPath)
	return err == nil
}

// Close closes the indexer (no-op for BleveIndexer as we don't keep indexes open)
func (b *BleveIndexer) Close() error {
	return nil
}

// getIndexPath returns the path to the index directory for a schema
func (b *BleveIndexer) getIndexPath(schemaID string) string {
	return filepath.Join(b.baseDir, schemaID+".bleve")
}

// buildIndexMapping creates the index mapping for schema documents
func buildIndexMapping() *mapping.IndexMappingImpl {
	// Create a text field mapping with standard analyzer
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = "standard"

	// Create a keyword field mapping (exact match, no analysis)
	keywordFieldMapping := bleve.NewKeywordFieldMapping()

	// Create document mapping
	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("type", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("name", textFieldMapping)
	docMapping.AddFieldMappingsAt("description", textFieldMapping)
	docMapping.AddFieldMappingsAt("path", textFieldMapping)
	docMapping.AddFieldMappingsAt("schemaID", keywordFieldMapping)

	// Create index mapping
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping

	return indexMapping
}

// extractDocuments extracts searchable documents from a schema
func extractDocuments(schemaID string, schema *gql.GraphQLSchema) []document {
	docs := []document{}

	// Index Query fields
	for name, field := range schema.Query {
		docs = append(docs, document{
			Type:        "Field",
			Name:        name,
			Description: field.Description(),
			Path:        "Query." + name,
			SchemaID:    schemaID,
		})
	}

	// Index Mutation fields
	for name, field := range schema.Mutation {
		docs = append(docs, document{
			Type:        "Field",
			Name:        name,
			Description: field.Description(),
			Path:        "Mutation." + name,
			SchemaID:    schemaID,
		})
	}

	// Index Objects
	for name, obj := range schema.Object {
		docs = append(docs, document{
			Type:        "Object",
			Name:        name,
			Description: obj.Description(),
			Path:        name,
			SchemaID:    schemaID,
		})
		docs = append(docs, extractObjectFieldDocuments(schemaID, name, obj)...)
	}

	// Index Input Objects
	for name, input := range schema.Input {
		docs = append(docs, document{
			Type:        "InputObject",
			Name:        name,
			Description: input.Description(),
			Path:        name,
			SchemaID:    schemaID,
		})
		docs = append(docs, extractInputFieldDocuments(schemaID, name, input)...)
	}

	// Index Enums
	for name, enum := range schema.Enum {
		docs = append(docs, document{
			Type:        "Enum",
			Name:        name,
			Description: enum.Description(),
			Path:        name,
			SchemaID:    schemaID,
		})
	}

	// Index Scalars
	for name, scalar := range schema.Scalar {
		docs = append(docs, document{
			Type:        "Scalar",
			Name:        name,
			Description: scalar.Description(),
			Path:        name,
			SchemaID:    schemaID,
		})
	}

	// Index Interfaces
	for name, iface := range schema.Interface {
		docs = append(docs, document{
			Type:        "Interface",
			Name:        name,
			Description: iface.Description(),
			Path:        name,
			SchemaID:    schemaID,
		})
		docs = append(docs, extractInterfaceFieldDocuments(schemaID, name, iface)...)
	}

	// Index Unions
	for name, union := range schema.Union {
		docs = append(docs, document{
			Type:        "Union",
			Name:        name,
			Description: union.Description(),
			Path:        name,
			SchemaID:    schemaID,
		})
	}

	// Index Directives
	for name, directive := range schema.Directive {
		docs = append(docs, document{
			Type:        "Directive",
			Name:        name,
			Description: directive.Description(),
			Path:        "@" + name,
			SchemaID:    schemaID,
		})
	}

	return docs
}

// extractObjectFieldDocuments extracts documents for fields of an object
func extractObjectFieldDocuments(schemaID string, typeName string, obj *gql.Object) []document {
	docs := []document{}
	for _, field := range obj.Fields() {
		docs = append(docs, document{
			Type:        "Field",
			Name:        field.Name(),
			Description: field.Description(),
			Path:        fmt.Sprintf("%s.%s", typeName, field.Name()),
			SchemaID:    schemaID,
		})
	}
	return docs
}

// extractInputFieldDocuments extracts documents for fields of an input object
func extractInputFieldDocuments(schemaID string, typeName string, input *gql.InputObject) []document {
	docs := []document{}
	for _, field := range input.Fields() {
		docs = append(docs, document{
			Type:        "Field",
			Name:        field.Name(),
			Description: field.Description(),
			Path:        fmt.Sprintf("%s.%s", typeName, field.Name()),
			SchemaID:    schemaID,
		})
	}
	return docs
}

// extractInterfaceFieldDocuments extracts documents for fields of an interface
func extractInterfaceFieldDocuments(schemaID string, typeName string, iface *gql.Interface) []document {
	docs := []document{}
	for _, field := range iface.Fields() {
		docs = append(docs, document{
			Type:        "Field",
			Name:        field.Name(),
			Description: field.Description(),
			Path:        fmt.Sprintf("%s.%s", typeName, field.Name()),
			SchemaID:    schemaID,
		})
	}
	return docs
}
