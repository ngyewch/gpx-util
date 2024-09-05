package main

import (
	gpxutil "github.com/ngyewch/gpx-util"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"time"
)

var (
	convertCmd = &cli.Command{
		Name:      "convert",
		Usage:     "Convert GPS NMEA to GPX",
		ArgsUsage: "[(input file) [(output file)]]",
		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name:     "start-date",
				Usage:    "start date (for first point, UTC)",
				Layout:   time.DateOnly,
				Timezone: time.UTC,
			},
		},
		Action: doConvert,
	}
)

func doConvert(cCtx *cli.Context) error {
	startDate := cCtx.Timestamp("start-date")
	if startDate == nil {
		t := time.Now().UTC()
		startDate = &t
	}

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

	var w io.Writer
	if (cCtx.NArg() < 2) || (cCtx.Args().Get(1) == "-") {
		w = os.Stdout
	} else {
		f, err := os.Create(cCtx.Args().Get(1))
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		w = f
	}

	return gpxutil.Convert(*startDate, r, w)
}
