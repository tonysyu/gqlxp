package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tonysyu/gq/gql"
	"github.com/tonysyu/gq/tui"
)

func main() {
	// Read schema file (assuming it's in the project root or data directory)
	schemaContent, err := ioutil.ReadFile("examples/github.graphqls")
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	schema := gql.ParseSchema(schemaContent)

	if _, err := tui.Start(schema); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
