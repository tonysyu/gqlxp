package library

import "time"

// SchemaMetadata contains metadata for a stored schema.
type SchemaMetadata struct {
	DisplayName string            `json:"displayName"`
	SourceFile  string            `json:"sourceFile"`
	FileHash    string            `json:"fileHash"`
	URLPatterns map[string]string `json:"urlPatterns"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// Schema represents a stored schema with its content and metadata.
type Schema struct {
	ID       string
	Content  []byte
	Metadata SchemaMetadata
}

// SchemaInfo represents basic schema information for listing.
type SchemaInfo struct {
	ID          string
	DisplayName string
}

// UserConfig contains user preferences and settings.
type UserConfig struct {
	DefaultSchema string `json:"defaultSchema,omitempty"`
}
