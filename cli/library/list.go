package library

import (
	"context"
	"fmt"
	"strings"

	"github.com/tonysyu/gqlxp/library"
	"github.com/urfave/cli/v3"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all schemas in the library",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			lib := library.NewLibrary()
			schemas, err := lib.List()
			if err != nil {
				return fmt.Errorf("failed to list schemas: %w", err)
			}

			if len(schemas) == 0 {
				fmt.Println("No schemas in library. Add one with: gqlxp library add <schema-file>")
				return nil
			}

			// Get default schema to mark it
			defaultID, _ := lib.GetDefaultSchema()

			fmt.Println("Schemas in library:")
			for _, schema := range schemas {
				marker := " "
				if schema.ID == defaultID {
					marker = "*"
				}
				var parts []string
				if schema.DisplayName != "" {
					parts = append(parts, schema.DisplayName)
				}
				if !schema.UpdatedAt.IsZero() {
					parts = append(parts, "last-updated: "+schema.UpdatedAt.Format("2006-01-02 15:04"))
				}
				if len(parts) > 0 {
					fmt.Printf("%s %s (%s)\n", marker, schema.ID, strings.Join(parts, "; "))
				} else {
					fmt.Printf("%s %s\n", marker, schema.ID)
				}
			}

			return nil
		},
	}
}
