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
	Kind        string   `json:"kind"`        // Structural kind: Object, Query, ObjectField, etc.
	Name        string   `json:"name"`        // Type or field name
	Description string   `json:"description"` // Description text
	Path        string   `json:"path"`        // Full path (e.g., "Query.user.name")
	SchemaID    string   `json:"schemaID"`    // Schema identifier
	Signature   string   `json:"signature"`   // Field signature (e.g., "getUser(id: ID!): User")
	Usage       []string `json:"usage"`       // ParentKinds of types/fields that reference this type
	Implements  []string `json:"implements"`  // Interface names this type implements (Object/Interface only)
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
	docMapping.AddFieldMappingsAt("kind", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("name", textFieldMapping)
	docMapping.AddFieldMappingsAt("description", textFieldMapping)
	docMapping.AddFieldMappingsAt("path", textFieldMapping)
	docMapping.AddFieldMappingsAt("schemaID", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("signature", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("usage", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("implements", keywordFieldMapping)

	// Create index mapping
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping

	return indexMapping
}

// extractDocuments extracts searchable documents from a schema
func extractDocuments(schemaID string, schema *gql.GraphQLSchema) []document {
	var docs []document

	schema.Walk(gql.SchemaVisitor{
		VisitField: func(ctx gql.VisitContext, name string, field *gql.Field) {
			docs = append(docs, document{
				Kind:        ctx.Kind,
				Name:        name,
				Description: field.Description(),
				Path:        ctx.Kind + "." + name,
				SchemaID:    schemaID,
				Signature:   field.Signature(),
				Usage:       fieldUsage(field),
			})
		},
		VisitObject: func(ctx gql.VisitContext, name string, obj *gql.Object) {
			docs = append(docs, document{
				Kind:        "Object",
				Name:        name,
				Description: obj.Description(),
				Path:        name,
				SchemaID:    schemaID,
				Implements:  obj.Interfaces(),
			})
		},
		VisitObjectField: func(ctx gql.VisitContext, field *gql.Field) {
			docs = append(docs, document{
				Kind:        "ObjectField",
				Name:        field.Name(),
				Description: field.Description(),
				Path:        ctx.ParentName + "." + field.Name(),
				SchemaID:    schemaID,
				Signature:   field.Signature(),
				Usage:       fieldUsage(field),
			})
		},
		VisitInterface: func(ctx gql.VisitContext, name string, iface *gql.Interface) {
			docs = append(docs, document{
				Kind:        "Interface",
				Name:        name,
				Description: iface.Description(),
				Path:        name,
				SchemaID:    schemaID,
				Implements:  iface.Interfaces(),
			})
		},
		VisitInterfaceField: func(ctx gql.VisitContext, field *gql.Field) {
			docs = append(docs, document{
				Kind:        "InterfaceField",
				Name:        field.Name(),
				Description: field.Description(),
				Path:        ctx.ParentName + "." + field.Name(),
				SchemaID:    schemaID,
				Signature:   field.Signature(),
				Usage:       fieldUsage(field),
			})
		},
		VisitInput: func(ctx gql.VisitContext, name string, input *gql.InputObject) {
			docs = append(docs, document{
				Kind:        "Input",
				Name:        name,
				Description: input.Description(),
				Path:        name,
				SchemaID:    schemaID,
			})
		},
		VisitInputField: func(ctx gql.VisitContext, field *gql.Field) {
			docs = append(docs, document{
				Kind:        "InputField",
				Name:        field.Name(),
				Description: field.Description(),
				Path:        ctx.ParentName + "." + field.Name(),
				SchemaID:    schemaID,
				Signature:   field.Signature(),
				Usage:       fieldUsage(field),
			})
		},
		VisitEnum: func(ctx gql.VisitContext, name string, enum *gql.Enum) {
			docs = append(docs, document{
				Kind:        "Enum",
				Name:        name,
				Description: enum.Description(),
				Path:        name,
				SchemaID:    schemaID,
			})
		},
		VisitScalar: func(ctx gql.VisitContext, name string, scalar *gql.Scalar) {
			docs = append(docs, document{
				Kind:        "Scalar",
				Name:        name,
				Description: scalar.Description(),
				Path:        name,
				SchemaID:    schemaID,
			})
		},
		VisitUnion: func(ctx gql.VisitContext, name string, union *gql.Union) {
			docs = append(docs, document{
				Kind:        "Union",
				Name:        name,
				Description: union.Description(),
				Path:        name,
				SchemaID:    schemaID,
			})
		},
		VisitDirective: func(ctx gql.VisitContext, name string, directive *gql.DirectiveDef) {
			docs = append(docs, document{
				Kind:        "Directive",
				Name:        name,
				Description: directive.Description(),
				Path:        "@" + name,
				SchemaID:    schemaID,
			})
		},
	})

	return docs
}

// fieldUsage returns a slice containing the field's return type name, or nil if the type is a
// built-in scalar (no object type name).
func fieldUsage(field *gql.Field) []string {
	typeName := field.ObjectTypeName()
	if typeName == "" {
		return nil
	}
	return []string{typeName}
}
