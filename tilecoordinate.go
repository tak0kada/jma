package jma

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"math"
)

type TileCoordinate struct {
	Zoom uint    `validate:"gte=0,lte=18,required"`
	X    float64 `json:"tilecoordinate_x" validate:"required"`
	Y    float64 `json:"tilecoordinate_y" validate:"required"`
}

func (tc TileCoordinate) String() string {
	return fmt.Sprintf("{Level: %d, X: %g, Y: %g}", tc.Zoom, tc.X, tc.Y)
}

func (tc TileCoordinate) GetTile() Tile {
	return Tile{
		Zoom: tc.Zoom,
		X:    uint(tc.X),
		Y:    uint(tc.Y),
	}
}

func (tc TileCoordinate) ToGeoCoordinate() GeoCoordinate {
	n := math.Pow(2, float64(tc.Zoom))
	my := 2*tc.Y*math.Pi/n - math.Pi
	return GeoCoordinate{
		Lat: math.Atan(math.Exp(-my))*360/math.Pi - 90,
		Lon: (tc.X/n)*360 - 180,
	}
}

func (tc TileCoordinate) IsValid() bool {
	err := validate.Struct(tc)
	if err != nil {
		// fmt.Printf("Err(s):\n%+v\n", err)
		return false
	}
	return true
}

func tileCoordinateXYValidation(sl validator.StructLevel) {
	tc := sl.Current().Interface().(TileCoordinate)
	if tc.X >= math.Pow(2, float64(tc.Zoom)) {
		sl.ReportError(tc.X, "tilecoordinate_x", "X", "x_should_be_less_than_2^level", "")
	}
	if tc.Y >= math.Pow(2, float64(tc.Zoom)) {
		sl.ReportError(tc.Y, "tilecoordinate_y", "Y", "y_should_be_less_than_2^level", "")
	}
}
