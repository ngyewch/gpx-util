package gpx_util

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/twpayne/go-gpx"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"io"
	"os"
)

func View(r io.Reader) error {
	g, err := gpx.Read(r)
	if err != nil {
		return err
	}

	var distanceXYs plotter.XYs
	var elevationXYs plotter.XYs
	var speedXYs plotter.XYs
	var courseXYs plotter.XYs
	var satellitesXYs plotter.XYs
	var hdopXYs plotter.XYs

	for _, trk := range g.Trk {
		for _, trkSeg := range trk.TrkSeg {
			distance := 0.0
			var p0 orb.Point
			var t0 int64
			for trkPtNo, trkPt := range trkSeg.TrkPt {
				if trkPtNo == 0 {
					t0 = trkPt.Time.Unix()
					p0 = orb.Point{trkPt.Lon, trkPt.Lat}
				}
				t := float64(trkPt.Time.Unix() - t0)

				p := orb.Point{trkPt.Lon, trkPt.Lat}
				dDistance := geo.DistanceHaversine(p0, p)
				distance += dDistance
				distanceXYs = append(distanceXYs, plotter.XY{
					X: t,
					Y: distance,
				})
				elevationXYs = append(elevationXYs, plotter.XY{
					X: t,
					Y: trkPt.Ele,
				})
				speedXYs = append(speedXYs, plotter.XY{
					X: t,
					Y: trkPt.Speed,
				})
				courseXYs = append(courseXYs, plotter.XY{
					X: t,
					Y: trkPt.Course,
				})
				satellitesXYs = append(satellitesXYs, plotter.XY{
					X: t,
					Y: float64(trkPt.Sat),
				})
				hdopXYs = append(hdopXYs, plotter.XY{
					X: t,
					Y: trkPt.HDOP,
				})

				p0 = p
			}
		}
		break
	}

	distanceChart, err := plotter.NewLine(distanceXYs)
	if err != nil {
		return err
	}
	elevationChart, err := plotter.NewLine(elevationXYs)
	if err != nil {
		return err
	}
	speedChart, err := plotter.NewLine(speedXYs)
	if err != nil {
		return err
	}
	courseChart, err := plotter.NewLine(courseXYs)
	if err != nil {
		return err
	}
	satellitesChart, err := plotter.NewLine(satellitesXYs)
	if err != nil {
		return err
	}
	hdopChart, err := plotter.NewLine(hdopXYs)
	if err != nil {
		return err
	}
	plots := []plot.Plotter{
		distanceChart,
		elevationChart,
		speedChart,
		courseChart,
		satellitesChart,
		hdopChart,
	}

	const chartWidth = 6 * vg.Inch
	const chartHeight = 4 * vg.Inch
	const rows = 6
	const cols = 1

	plotsTable := make([][]*plot.Plot, rows)
	for j := 0; j < rows; j++ {
		p := plot.New()
		p.X.Label.Text = "Time (s)"
		p.X.Min = 0
		//p.Y.Label.Text = "Total distance travelled (m)"
		p.Add(plots[j])

		plotsTable[j] = make([]*plot.Plot, cols)
		plotsTable[j][0] = p
	}

	img := vgimg.New(vg.Length(cols)*chartWidth, vg.Length(rows)*chartHeight)
	dc := draw.New(img)
	t := draw.Tiles{
		Rows:      rows,
		Cols:      cols,
		PadX:      1 * vg.Millimeter,
		PadY:      1 * vg.Millimeter,
		PadTop:    vg.Points(2),
		PadBottom: vg.Points(2),
		PadLeft:   vg.Points(2),
		PadRight:  vg.Points(2),
	}
	canvases := plot.Align(plotsTable, t, dc)
	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			if plotsTable[j][i] != nil {
				plotsTable[j][i].Draw(canvases[j][i])
			}
		}
	}

	f, err := os.Create("test.png")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	png := vgimg.PngCanvas{Canvas: img}
	_, err = png.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}
