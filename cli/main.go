package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	initcmd "github.com/tonysyu/gqlxp/cli/init"
	"github.com/tonysyu/gqlxp/cli/library"
	"github.com/tonysyu/gqlxp/tui"
)

// NewRootCmd creates and configures the CLI application.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gqlxp",
		Short: "Explore GraphQL schemas interactively or via CLI",
		Long: `gqlxp helps you explore, search, and validate GraphQL schemas.

For AI/programmatic use (use --ai flag for JSON output, no pager, no color):
  search    Find types and fields by keyword
  show      Display a full type definition
  validate  Validate a GraphQL operation against the schema
  generate  Scaffold a skeleton GraphQL operation (prints to stdout)

Schema files are saved to the library on first use.
Use 'gqlxp library list' to see available schemas.
Use 'gqlxp library default <id>' to set the default schema for all commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeTUICommand(cmd)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.PersistentFlags().StringP("schema", "s", "", "Schema file path or library ID")
	root.Flags().StringP("log-file", "l", "", "Enable debug logging to `FILE`")
	root.Flags().String("select", "", "Pre-select TYPE or TYPE.FIELD in TUI")

	root.AddCommand(
		appCommand(),
		initcmd.Command(),
		validateCommand(),
		searchCommand(),
		showCommand(),
		generateCommand(),
		library.Command(),
	)

	return root
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
