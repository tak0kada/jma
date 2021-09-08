package jma

import (
	"math"
	"testing"
)

func TestTileCoordinateToGeoCoordinate(t *testing.T) {
	tests := []struct {
		name string
		geoc GeoCoordinate
		zoom uint
	}{
		{
			name: "Akashi2/3/1",
			geoc: GeoCoordinate{34.6497512427944, 135.00132061562732},
			zoom: 2,
		},
		{
			name: "Akashi14/14336/6509",
			geoc: GeoCoordinate{34.6497512427944, 135.00132061562732},
			zoom: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.geoc.GetTileCoordinate(tt.zoom).ToGeoCoordinate(); !isApproxEq(tt.geoc, got) {
				t.Errorf("Geocoordinate=%s, expected=%s", tt.geoc, got)
			}
		})
	}
}

func isApproxEq(left GeoCoordinate, right GeoCoordinate) bool {
	return math.Abs(left.Lat-right.Lat)+math.Abs(left.Lon-right.Lon) < 1e-10
}
