package main

import (
	"io/ioutil"
	"log"

	"github.com/tonysyu/gq/gql"
)

func main() {
	// Read the GraphQL schema file
	schemaContent, err := ioutil.ReadFile("examples/github.graphqls")
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	// Parse and display the schema
	gql.ParseSchema(schemaContent)
}
