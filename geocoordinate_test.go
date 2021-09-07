package jma

import (
	"testing"
)

func TestGeocoordinateToTile(t *testing.T) {
	tests := []struct {
		name string
		geoc GeoCoordinate
		zoom uint
		tile Tile
	}{
		{
			name: "Akashi2/3/1",
			geoc: GeoCoordinate{34.6497512427944, 135.00132061562732},
			zoom: 2,
			tile: Tile{2, 3, 1},
		},
		{
			name: "Akashi14/14336/6509",
			geoc: GeoCoordinate{34.6497512427944, 135.00132061562732},
			zoom: 14,
			tile: Tile{14, 14336, 6509},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.geoc.ToTile(tt.zoom); got != tt.tile {
				t.Errorf("Geocoordinate=%s, Tile=%s, expected=%s", tt.geoc, got, tt.tile)
			}
		})
	}
}
