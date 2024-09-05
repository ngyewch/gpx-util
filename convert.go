package gpx_util

import (
	"bufio"
	"github.com/adrianmo/go-nmea"
	"github.com/twpayne/go-gpx"
	"io"
	"time"
)

func Convert(date time.Time, r io.Reader, w io.Writer) error {
	gpxFilter := NewGPXFilter(date)
	var trackPoints []*gpx.WptType
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		s, err := nmea.Parse(line)
		if err != nil {
			continue
		}
		point := gpxFilter.Handle(s)
		if point != nil {
			trackPoints = append(trackPoints, point)
		}
	}
	g := &gpx.GPX{
		Version: "1.0",
		Creator: "spd",
		Trk: []*gpx.TrkType{
			{
				TrkSeg: []*gpx.TrkSegType{
					{
						TrkPt: trackPoints,
					},
				},
			},
		},
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	err := g.WriteIndent(w, "", "  ")
	if err != nil {
		return err
	}
	return nil
}
