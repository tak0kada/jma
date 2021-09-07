package jma

import (
	"errors"
	"testing"
	"time"
)

func TestTileIsValid(t *testing.T) {
	tests := []struct {
		name string
		tile Tile
		want bool
	}{
		{
			name: "ok",
			tile: Tile{2, 3, 1},
			want: true,
		},
		{
			name: "ng",
			tile: Tile{2, 3, 1000},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tile.IsValid(); got != tt.want {
				t.Errorf("Tile=%s, result=%t, expected=%t", tt.tile, got, tt.want)
			}
		})
	}
}

func TestTileToMapURL(t *testing.T) {
	test := struct {
		name string
		tile Tile
		want string
	}{
		name: "test",
		tile: Tile{14, 14336, 6509},
		want: "https://www.jma.go.jp/tile/gsi/pale/14/14336/6509.png",
	}
	t.Run(test.name, func(t *testing.T) {
		if got := test.tile.ToMapURL("pale", "png"); got != test.want {
			t.Errorf("Tile=%s, result=%s, expected=%s", test.tile, got, test.want)
		}
	})
}

func TestTileToBorderMapURL(t *testing.T) {
	test := struct {
		name string
		tile Tile
		want string
	}{
		name: "test",
		tile: Tile{14, 14336, 6509},
		want: "https://www.jma.go.jp/bosai/jmatile/data/map/none/none/none/surf/mask/14/14336/6509.png",
	}
	t.Run(test.name, func(t *testing.T) {
		if got := test.tile.ToBorderMapURL("png"); got != test.want {
			t.Errorf("Tile=%s, result=%s, expected=%s", test.tile, got, test.want)
		}
	})
}

func TestTileToJmaURL(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2021-09-05T13:32:38Z")
	tests := []struct {
		name     string
		tile     Tile
		now      time.Time
		duration time.Duration
		want     string
		err      error
	}{
		{
			name:     "ok",
			tile:     Tile{14, 14336, 6509},
			now:      now,
			duration: -time.Hour,
			want:     "https://www.jma.go.jp/bosai/jmatile/data/nowc/20210905123000/none/20210905123000/surf/hrpns/14/14336/6509.png",
			err:      nil,
		},
		{
			name:     "ok",
			tile:     Tile{14, 14336, 6509},
			now:      now,
			duration: 30 * time.Minute,
			want:     "https://www.jma.go.jp/bosai/jmatile/data/nowc/20210905133000/none/20210905140000/surf/hrpns/14/14336/6509.png",
			err:      nil,
		},
		{
			name:     "ok",
			tile:     Tile{14, 14336, 6509},
			now:      now,
			duration: -4 * time.Hour,
			want:     "https://www.jma.go.jp/bosai/jmatile/data/rasrf/20210905090000/none/20210905090000/surf/rasrf/14/14336/6509.png",
			err:      nil,
		},
		{
			name:     "ok",
			tile:     Tile{14, 14336, 6509},
			now:      now,
			duration: 2 * time.Hour,
			want:     "https://www.jma.go.jp/bosai/jmatile/data/rasrf/20210905130000/none/20210905150000/surf/rasrf/14/14336/6509.png",
			err:      nil,
		},
		{
			name:     "ng",
			tile:     Tile{14, 14336, 6509},
			now:      now,
			duration: -24 * time.Hour,
			want:     "",
			err:      errors.New("forecasting is supported from -12hr to 14hr of duration only"),
		},
		{
			name:     "ng",
			tile:     Tile{14, 14336, 6509},
			now:      now,
			duration: 24 * time.Hour,
			want:     "",
			err:      errors.New("forecasting is supported from -12hr to 14hr of duration only"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.tile.ToJmaURL(tt.now, tt.duration, "png"); got != tt.want || err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Tile=%s, result=(%s, %s), expected=(%s, %s)", tt.tile, got, err, tt.want, tt.err)
			}
		})
	}
}
