package cli

import (
	"context"
	"fmt"
	"os"

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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-file",
				Aliases: []string{"l"},
				Usage:   "Enable debug logging to `FILE`",
				Sources: cli.EnvVars("GQLXP_LOGFILE"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return executeTUICommand(ctx, cmd)
		},
		Commands: []*cli.Command{
			appCommand(),
			searchCommand(),
			showCommand(),
			configCommand(),
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
