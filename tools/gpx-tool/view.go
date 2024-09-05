package main

import (
	gpx_util "github.com/ngyewch/gpx-util"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

var (
	viewCmd = &cli.Command{
		Name:      "view",
		Usage:     "View GPX file",
		ArgsUsage: "[(input file)]",
		Action:    doView,
	}
)

func doView(cCtx *cli.Context) error {
	var r io.Reader
	if (cCtx.NArg() < 1) || (cCtx.Args().Get(0) == "-") {
		r = os.Stdin
	} else {
		f, err := os.Open(cCtx.Args().Get(0))
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		r = f
	}

	return gpx_util.View(r)
}
