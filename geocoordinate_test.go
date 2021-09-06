package jma

import (
	"testing"
)

func TestGeocoordinateToTile(t *testing.T) {
	tests := []struct {
		name      string
		geoc      GeoCoordinate
		zoomLevel uint
		tile      Tile
	}{
		{
			name:      "Akashi2/3/1",
			geoc:      GeoCoordinate{34.6497512427944, 135.00132061562732},
			zoomLevel: 2,
			tile:      Tile{2, 3, 1},
		},
		{
			name:      "Akashi14/14336/6509",
			geoc:      GeoCoordinate{34.6497512427944, 135.00132061562732},
			zoomLevel: 14,
			tile:      Tile{14, 14336, 6509},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.geoc.ToTile(tt.zoomLevel); got != tt.tile {
				t.Errorf("Geocoordinate=%s, Tile=%s, expected=%s", tt.geoc.String(), got.String(), tt.tile.String())
			}
		})
	}
}
