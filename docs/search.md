# Schema Search

Fast full-text search across GraphQL schema types and fields.

## Usage

```sh
# Search in a schema file
gqlxp search -s examples/github.graphqls user

# Search using schema from library
gqlxp search -s github-api user

# Search using default schema (omit --schema flag)
gqlxp search mutation

# Options
gqlxp search --limit 5 user    # Limit results (default: 30)
```

Results show type, name, path, and description ranked by relevance.

## Indexing

Schemas are automatically indexed when added or updated. Indexes are stored in `~/.config/gqlxp/schemas/<schema-id>.bleve/` and removed when schemas are deleted.

To rebuild an index manually:
```sh
gqlxp library reindex <schema-id>
gqlxp library reindex --all  # Rebuild all indexes
```

Search covers all schema elements: queries, mutations, objects, input objects, interfaces, enums, scalars, unions, and directives.
