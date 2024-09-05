package gpx_util

import (
	"github.com/adrianmo/go-nmea"
	"github.com/martinlindhe/unit"
	"github.com/twpayne/go-gpx"
	"strconv"
	"time"
)

type GPXFilter struct {
	date     time.Time
	lastTime time.Time
	gga      *nmea.GGA
	rmc      *nmea.RMC
}

func NewGPXFilter(startDate time.Time) *GPXFilter {
	d := startDate.UTC()
	return &GPXFilter{
		date: time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC),
	}
}

func (filter *GPXFilter) Handle(s nmea.Sentence) *gpx.WptType {
	{
		rmc, ok := s.(nmea.RMC)
		if ok {
			if (filter.rmc != nil) || ((filter.gga != nil) && (filter.gga.Time != rmc.Time)) {
				point := filter.flush()
				filter.rmc = &rmc
				return point
			} else {
				filter.rmc = &rmc
			}
		}
	}
	{
		gga, ok := s.(nmea.GGA)
		if ok {
			if (filter.gga != nil) || ((filter.rmc != nil) && (filter.rmc.Time != gga.Time)) {
				point := filter.flush()
				filter.gga = &gga
				return point
			} else {
				filter.gga = &gga
			}
		}
	}
	if (filter.gga != nil) && (filter.rmc != nil) {
		return filter.flush()
	}
	return nil
}

func (filter *GPXFilter) flush() *gpx.WptType {
	defer func() {
		filter.rmc = nil
		filter.gga = nil
	}()

	if (filter.rmc == nil) && (filter.gga == nil) {
		return nil
	}

	var point gpx.WptType
	if filter.rmc != nil {
		if filter.rmc.Validity != "A" {
			return nil
		}
		point.Lat = filter.rmc.Latitude
		point.Lon = filter.rmc.Longitude
		point.Speed = (unit.Speed(filter.rmc.Speed) * unit.Knot).MetersPerSecond()
		point.Course = filter.rmc.Course
		point.Time = filter.toTime(filter.rmc.Time)
		point.MagVar = filter.rmc.Variation
	}
	if filter.gga != nil {
		point.Lat = filter.gga.Latitude
		point.Lon = filter.gga.Longitude
		point.Time = filter.toTime(filter.gga.Time)
		switch filter.gga.FixQuality {
		case "0": // Fix not available or invalid
			return nil
		case "1": // AGPS SPS Mode, fix valid
			point.Fix = "3d"
		case "2": // Differential GPS, SPS Mode, fix valid
			point.Fix = "dgps"
		}
		point.Sat = int(filter.gga.NumSatellites)
		point.HDOP = filter.gga.HDOP
		point.Ele = filter.gga.Altitude
		point.GeoidHeight = filter.gga.Separation
		if filter.gga.DGPSAge != "" {
			dgpsAge, err := strconv.ParseFloat(filter.gga.DGPSAge, 64)
			if err == nil {
				point.AgeOfDGPSData = dgpsAge
			}
		}
		if filter.gga.DGPSId != "" {
			dgpsId, err := strconv.ParseInt(filter.gga.DGPSId, 10, 64)
			if err == nil {
				point.DGPSID = []int{int(dgpsId)}
			}
		}
	}

	return &point
}

func (filter *GPXFilter) toTime(nmeaTime nmea.Time) time.Time {
	t := toTime(filter.date, nmeaTime)
	if filter.lastTime.IsZero() || !t.Before(filter.lastTime) {
		filter.lastTime = t
		return t
	}
	filter.date = filter.date.AddDate(0, 0, 1)
	t = toTime(filter.date, nmeaTime)
	filter.lastTime = t
	return t
}

func toTime(date time.Time, nmeaTime nmea.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), nmeaTime.Hour, nmeaTime.Minute, nmeaTime.Second, nmeaTime.Millisecond*1_000_000, time.UTC)
}
