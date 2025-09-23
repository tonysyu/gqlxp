package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tonysyu/igq/gql"
	"github.com/tonysyu/igq/tui"
)

func main() {
	// Read schema file (assuming it's in the project root or data directory)
	schemaContent, err := ioutil.ReadFile("examples/github.graphqls")
	if err != nil {
		abort(fmt.Sprintf("Error reading schema file: %v\n", err))
	}

	schema, err := gql.ParseSchema(schemaContent)
	if err != nil {
		abort(fmt.Sprintf("Error parsing schema: %v", err))
	}

	if _, err := tui.Start(schema); err != nil {
		abort(fmt.Sprintf("Error starting tui: %v", err))
	}
}

func abort(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
