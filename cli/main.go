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
		Usage: "Explore GraphQL schemas interactively",
		Description: `Schema files are automatically saved to your library on first use.
When loading a previously imported file, you'll be prompted to update
if changes are detected.

Use the TUI interface to manage library schemas (remove, view, etc).`,
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
