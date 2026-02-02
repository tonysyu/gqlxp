package introspection

// Response represents the full introspection query response.
type Response struct {
	Data   *Data   `json:"data"`
	Errors []Error `json:"errors,omitempty"`
}

// Data wraps the schema in the response.
type Data struct {
	Schema Schema `json:"__schema"`
}

// Schema represents the introspected GraphQL schema.
type Schema struct {
	QueryType        *TypeName   `json:"queryType"`
	MutationType     *TypeName   `json:"mutationType"`
	SubscriptionType *TypeName   `json:"subscriptionType"`
	Types            []FullType  `json:"types"`
	Directives       []Directive `json:"directives"`
}

// TypeName is a simple type reference containing only the name.
type TypeName struct {
	Name string `json:"name"`
}

// FullType represents a complete type definition.
type FullType struct {
	Kind          string       `json:"kind"`
	Name          string       `json:"name"`
	Description   *string      `json:"description"`
	Fields        []Field      `json:"fields"`
	InputFields   []InputValue `json:"inputFields"`
	Interfaces    []TypeRef    `json:"interfaces"`
	EnumValues    []EnumValue  `json:"enumValues"`
	PossibleTypes []TypeRef    `json:"possibleTypes"`
}

// Field represents a field on an object or interface type.
type Field struct {
	Name              string       `json:"name"`
	Description       *string      `json:"description"`
	Args              []InputValue `json:"args"`
	Type              TypeRef      `json:"type"`
	IsDeprecated      bool         `json:"isDeprecated"`
	DeprecationReason *string      `json:"deprecationReason"`
}

// InputValue represents an input field or argument.
type InputValue struct {
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Type         TypeRef `json:"type"`
	DefaultValue *string `json:"defaultValue"`
}

// TypeRef represents a type reference with potential wrapping.
type TypeRef struct {
	Kind   string   `json:"kind"`
	Name   *string  `json:"name"`
	OfType *TypeRef `json:"ofType"`
}

// EnumValue represents a value in an enum type.
type EnumValue struct {
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	IsDeprecated      bool    `json:"isDeprecated"`
	DeprecationReason *string `json:"deprecationReason"`
}

// Directive represents a directive definition.
type Directive struct {
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	Locations   []string     `json:"locations"`
	Args        []InputValue `json:"args"`
}

// Error represents a GraphQL error in the response.
type Error struct {
	Message string `json:"message"`
}
