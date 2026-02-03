package library

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gql/introspection"
	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

// Command creates the library subcommand.
func Command() *cli.Command {
	return &cli.Command{
		Name:  "library",
		Usage: "Manage schema library",
		Description: `Centralized interface for managing the schema library.

Available subcommands:
  list     - List all schemas in the library
  add      - Add a schema to the library
  update   - Update a schema in the library
  remove   - Remove a schema from the library
  default  - Set or show the default schema
  reindex  - Rebuild search indexes`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Default action is to list schemas
			return listCommand().Run(ctx, []string{})
		},
		Commands: []*cli.Command{
			listCommand(),
			addCommand(),
			updateCommand(),
			removeCommand(),
			defaultCommand(),
			reindexCommand(),
		},
	}
}

// schemaSource represents the source of a schema (either URL or file path).
type schemaSource struct {
	URL      string
	FilePath string
}

// loadSchemaContent loads schema content from a file path or URL.
func loadSchemaContent(ctx context.Context, source string, headers []string) ([]byte, schemaSource, error) {
	if introspection.IsURL(source) {
		content, err := fetchSchemaFromURL(ctx, source, headers)
		if err != nil {
			return nil, schemaSource{}, err
		}
		return content, schemaSource{URL: source}, nil
	}

	// Load from file
	absPath, err := filepath.Abs(source)
	if err != nil {
		return nil, schemaSource{}, fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	content, err := loadSchemaFromFile(absPath)
	if err != nil {
		return nil, schemaSource{}, err
	}
	return content, schemaSource{FilePath: absPath}, nil
}

// fetchSchemaFromURL fetches a GraphQL schema via introspection from the given URL.
func fetchSchemaFromURL(ctx context.Context, endpoint string, headers []string) ([]byte, error) {
	opts := introspection.DefaultClientOptions()

	// Parse and add custom headers
	if len(headers) > 0 {
		customHeaders, err := introspection.ParseHeaders(headers)
		if err != nil {
			return nil, fmt.Errorf("failed to parse headers: %w", err)
		}
		for k, v := range customHeaders {
			opts.Headers[k] = v
		}
	}

	fmt.Printf("Fetching schema from %s...\n", endpoint)

	resp, err := introspection.FetchSchema(ctx, endpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schema: %w", err)
	}

	sdl, err := introspection.ToSDL(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to SDL: %w", err)
	}

	return sdl, nil
}

// getSchemaID returns schema ID from flag or prompts user.
func getSchemaID(cmd *cli.Command, source string) (string, error) {
	if flagID := cmd.String("id"); flagID != "" {
		if err := library.ValidateSchemaID(flagID); err != nil {
			return "", err
		}
		return flagID, nil
	}

	var suggested string
	if introspection.IsURL(source) {
		suggested = extractHostnameAsID(source)
	} else {
		basename := filepath.Base(source)
		ext := filepath.Ext(basename)
		suggested = strings.TrimSuffix(basename, ext)
	}
	suggested = library.SanitizeSchemaID(suggested)

	schemaID, err := promptSchemaID(suggested)
	if err != nil {
		return "", fmt.Errorf("failed to get schema ID: %w", err)
	}

	return schemaID, nil
}

// extractHostnameAsID extracts the hostname from a URL and sanitizes it for use as an ID.
func extractHostnameAsID(urlStr string) string {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "schema"
	}
	hostname := parsed.Hostname()
	hostname = strings.TrimPrefix(hostname, "api.")
	hostname = strings.TrimPrefix(hostname, "www.")
	parts := strings.Split(hostname, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return hostname
}

// getDisplayName returns display name from flag or prompts user.
func getDisplayName(cmd *cli.Command, defaultName string) (string, error) {
	if flagName := cmd.String("name"); flagName != "" {
		return flagName, nil
	}

	displayName, err := promptString("Enter display name", defaultName)
	if err != nil {
		return "", fmt.Errorf("failed to get display name: %w", err)
	}

	return displayName, nil
}

// schemaNotFoundError returns an error with the schema ID and lists available schemas.
func schemaNotFoundError(lib library.Library, schemaID string) error {
	schemas, err := lib.List()
	if err != nil || len(schemas) == 0 {
		return fmt.Errorf("schema '%s' not found in library", schemaID)
	}

	var ids []string
	for _, s := range schemas {
		ids = append(ids, s.ID)
	}
	return fmt.Errorf("schema '%s' not found in library. Available: %s", schemaID, strings.Join(ids, ", "))
}

// validateSchema parses content and returns error if invalid.
func validateSchema(content []byte) error {
	if _, err := gql.ParseSchema(content); err != nil {
		return fmt.Errorf("invalid GraphQL schema: %w", err)
	}
	return nil
}
