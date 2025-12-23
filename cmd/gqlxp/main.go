package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tonysyu/gqlxp/cli"
)

func main() {
	if err := cli.NewApp().Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
