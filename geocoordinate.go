package jma

import (
	"fmt"
	"math"
)

type GeoCoordinate struct {
	Lat float64 `validate:"gte=-85.0511287798,lte=85.0511287798,required"`
	Lon float64 `validate:"required"`
}

func (g GeoCoordinate) String() string {
	return fmt.Sprintf("{lat: %f, lon:%f}", g.Lat, g.Lon)
}

func (g GeoCoordinate) IsValid() bool {
	err := validate.Struct(g)
	if err != nil {
		// fmt.Printf("Err(s):\n%+v\n", err)
		return false
	}
	return true
}

func (g GeoCoordinate) Normalize() GeoCoordinate {
	latmax := 85.0511287798
	return GeoCoordinate{
		Lat: math.Max(math.Min(g.Lat, latmax), -latmax),
		Lon: math.Mod(g.Lon+180, 360) - 180,
	}
}

func (g GeoCoordinate) ToTile(zoom uint) Tile {
	// reference for conversion formula: https://standardization.at.webry.info/201401/article_1.html
	n := math.Pow(2, float64(zoom))
	r := math.Pi / 180
	lat := g.Lat * r
	return Tile{
		Zoom: zoom,
		X:    uint((g.Lon + 180) / 360 * n),
		Y:    uint((1.0 - math.Log(math.Tan(math.Pi/4+lat/2))/math.Pi) / 2 * n),
	}
}
