package main

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/tonysyu/gqlxp/cli"
)

func main() {
	if err := fang.Execute(context.Background(), cli.NewRootCmd()); err != nil {
		os.Exit(1)
	}
}
