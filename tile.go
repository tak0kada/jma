package jma

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"math"
	"time"
)

// For the definition of Tile struct, see: https://maps.gsi.go.jp/development/siyou.html#siyou-url .
type Tile struct {
	Zoom uint `validate:"gte=0,lte=18,required"`
	X    uint `json:"tile_x" validate:"required"`
	Y    uint `json:"tile_y" validate:"required"`
}

func (t Tile) String() string {
	return fmt.Sprintf("{Level: %d, X: %d, Y: %d}", t.Zoom, t.X, t.Y)
}

// Geospatial Information Authority(国土地理院)
func (t Tile) ToMapURL(datatype string, ext string) string {
	// For the valid pattern, see: https://maps.gsi.go.jp/development/ichiran.html.
	return fmt.Sprintf("https://www.jma.go.jp/tile/gsi/%s/%d/%d/%d."+ext,
		datatype, t.Zoom, t.X, t.Y)
}

func (t Tile) ToBorderMapURL(ext string) string {
	return fmt.Sprintf("https://www.jma.go.jp/bosai/jmatile/data/map/none/none/none/surf/mask/%d/%d/%d."+ext,
		t.Zoom, t.X, t.Y)
}

// Japan Meteorological Agency(気象庁)
func (t Tile) ToJmaURL(now time.Time, duration time.Duration, ext string) (string, error) {
	if -3*time.Hour < duration && duration < time.Hour {
		return t.toJmaHighresURL(now, duration, ext), nil
	} else if -12*time.Hour < duration && duration < 15*time.Hour {
		return t.toJmaLowresURL(now, duration, ext), nil
	} else {
		return "", errors.New("forecasting is supported from -12hr to 14hr of duration only")
	}
}

func formatTime(time time.Time) string {
	return time.UTC().Format("20060102150405")
}

func (t Tile) toJmaLowresURL(now time.Time, duration time.Duration, ext string) string {
	var after time.Time
	round := now.Round(time.Hour)
	if round.Before(now) {
		now = round
		after = now.Add(duration)
	} else {
		now = round.Add(-1 * time.Hour)
		after = now.Add(duration)
	}
	if duration <= 0 {
		return fmt.Sprintf("https://www.jma.go.jp/bosai/jmatile/data/rasrf/%s/none/%s/surf/rasrf/%d/%d/%d.png",
			formatTime(after), formatTime(after), t.Zoom, t.X, t.Y)
	} else {
		return fmt.Sprintf("https://www.jma.go.jp/bosai/jmatile/data/rasrf/%s/none/%s/surf/rasrf/%d/%d/%d.png",
			formatTime(now), formatTime(after), t.Zoom, t.X, t.Y)
	}
}

func (t Tile) toJmaHighresURL(now time.Time, duration time.Duration, ext string) string {
	var after time.Time
	round := now.Round(5 * time.Minute)
	if round.Before(now) {
		now = round
		after = now.Add(duration)
	} else {
		now = round.Add(-5 * time.Minute)
		after = now.Add(duration)
	}
	if duration <= 0 {
		return fmt.Sprintf("https://www.jma.go.jp/bosai/jmatile/data/nowc/%s/none/%s/surf/hrpns/%d/%d/%d.png",
			formatTime(after), formatTime(after), t.Zoom, t.X, t.Y)
	} else {
		return fmt.Sprintf("https://www.jma.go.jp/bosai/jmatile/data/nowc/%s/none/%s/surf/hrpns/%d/%d/%d.png",
			formatTime(now), formatTime(after), t.Zoom, t.X, t.Y)
	}
}

func (t Tile) IsValid() bool {
	err := validate.Struct(t)
	if err != nil {
		// fmt.Printf("Err(s):\n%+v\n", err)
		return false
	}
	return true
}

func tileXYValidation(sl validator.StructLevel) {
	tile := sl.Current().Interface().(Tile)
	if tile.X >= uint(math.Pow(2, float64(tile.Zoom))) {
		sl.ReportError(tile.X, "tile_x", "X", "x_should_be_less_than_2^level", "")
	}
	if tile.Y >= uint(math.Pow(2, float64(tile.Zoom))) {
		sl.ReportError(tile.Y, "tile_y", "Y", "y_should_be_less_than_2^level", "")
	}
}
