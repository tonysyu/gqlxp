package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/tonysyu/gqlxp/cli"
)

func main() {
	if err := fang.Execute(context.Background(), cli.NewRootCmd()); err != nil {
		os.Exit(1)
	}
}
