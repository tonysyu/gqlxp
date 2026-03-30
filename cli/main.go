package cli

import (
	"context"
	"fmt"
	"os"

	initcmd "github.com/tonysyu/gqlxp/cli/init"
	"github.com/tonysyu/gqlxp/cli/library"
	"github.com/tonysyu/gqlxp/tui"
	"github.com/urfave/cli/v3"
)

// NewApp creates and configures the CLI application.
func NewApp() *cli.Command {
	return &cli.Command{
		Name:  "gqlxp",
		Usage: "Explore GraphQL schemas interactively or via CLI",
		Description: `gqlxp helps you explore, search, and validate GraphQL schemas.

For AI/programmatic use:
  search    Find types and fields by keyword (--json --no-pager for structured output)
  show      Display a full type definition (--json --no-pager for structured output)
  validate  Validate a GraphQL operation against the schema (--json for structured output)
  generate  Scaffold a skeleton GraphQL operation (prints to stdout)

Schema files are saved to the library on first use.
Use 'gqlxp library list' to see available schemas.
Use 'gqlxp library default <id>' to set the default schema for all commands.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return executeTUICommand(ctx, cmd)
		},
		Commands: []*cli.Command{
			appCommand(),
			initcmd.Command(),
			validateCommand(),
			searchCommand(),
			showCommand(),
			generateCommand(),
			library.Command(),
		},
	}
}

func setupLogging(logFile string) {
	err := tui.SetupLogging(logFile)
	if err != nil {
		abort(fmt.Sprintf("Error opening log file: %v", err))
	}
}

func abort(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
