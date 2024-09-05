package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"runtime/debug"
)

func main() {
	app := &cli.App{
		Name:  "gpx-util",
		Usage: "GPX util",
		Commands: []*cli.Command{
			convertCmd,
			viewCmd,
		},
	}

	buildInfo, _ := debug.ReadBuildInfo()
	if buildInfo != nil {
		app.Version = buildInfo.Main.Version
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
