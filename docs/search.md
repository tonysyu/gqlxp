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

## Search Syntax

Search is implemented using [bleve](https://github.com/blevesearch/bleve) and supports `bleve`'s [query syntax](https://blevesearch.com/docs/Query-String-Query/).

The following examples assume a default schema has been set using `gqlxp library default`:
```sh
# Search for Query related to "user"
$ gqlxp search "type:Query user"

# Search for Query with name "user" (results match _either_ type or name)
$ gqlxp search "type:Query name:user"

# Search for Query with name "user" (results match _both_ type and name)
$ gqlxp search "+type:Query +name:user"

# Search for Mutation with name _containing_ "user"
$ gqlxp search "+type:Mutation +name:*user*"

# Search for Object or Interface with name starting with "repo" (slashes denote regex)
$ gqlxp search "+type:/(Object|Interface)/ +name:repo*"

# Search for "user" in all GraphQL types and exclude all fields
$ gqlxp search "user -type:*Field"
```

### Searching Word Fragments

By default, bleve does not search word fragments. The search syntax examples above show examples using a wild-card query (e.g. `*`) to search for word fragments. (Something similar can be done w/ regexes using `.*`.)

Since searching for fragments of type names is common, logic in gqlxp special cases
"simple searches" to search for fragments in the `name` field. For example, searching
for `myterm` is roughly equivalent to searching for `myterm name:*myterm*`.

This default word-fragment behavior is only applied to "simple" searches, which
requires no whitespace is used and that none of the following characters are in the query `:`, `*`, `?`, `"`, `+`, `-`, `(`, `)`, `~`, `^`, '\\'.

### Search fields

Each search document contains the following fields:
- `type`: GraphQL type or field (see section below)
- `name`: Name of the type or field
- `description`: Description/docstring of type or field
- `path`: Qualified name, which matches `name` for types and `<parent-name>.<name>` for fields
    - Queries and mutations will have fixed paths `Query.<name>` and `Mutation.<name>`,
      respectively

### Document `type`s

You can specify the standard GraphQL types for the `type` field in your query:
- `Query`
- `Mutation`
- `Object`
- `Input`
- `Enum`
- `Scalar`
- `Interface`
- `Union`
- `Directive`

In addition to the standard types, the following field types are defined:
- `ObjectField`
- `InputField`
- `InterfaceField`
