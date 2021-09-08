package jma

import (
	"github.com/tak0kada/jma/internal"
	"image"
	"math"
	"time"
)

func FetchImage(gc GeoCoordinate, zoom uint, rect Rect, now time.Time, duration time.Duration) (image.Image, error) {
	base, err := FetchMapImage(gc, zoom, rect, "pale")
	if err != nil {
		return nil, err
	}
	weather, err := FetchJmaImage(gc, zoom, rect, now, duration)
	if err != nil {
		return nil, err
	}
	border, err := FetchBorderImage(gc, zoom, rect)
	if err != nil {
		return nil, err
	}
	img, _ := Overlay(background(base), weather, border)
	return img, nil
}

func FetchImageTile(tile Tile, now time.Time, duration time.Duration) (image.Image, error) {
	tc := TileCoordinate{
		Zoom: tile.Zoom,
		X:    float64(tile.X) + 0.5,
		Y:    float64(tile.Y) + 0.5,
	}
	return FetchImage(tc.ToGeoCoordinate(), tc.Zoom, Rect{256, 256}, now, duration)
}

func FetchMapImage(gc GeoCoordinate, zoom uint, rect Rect, datatype string) (image.Image, error) {
	tc := gc.GetTileCoordinate(zoom)
	tiles := makeTiles(tc, rect)
	imgs := make([][]image.Image, len(tiles))
	for h := range imgs {
		imgs[h] = make([]image.Image, len(tiles[0]))
	}
	for h := range imgs {
		for w := range imgs[h] {
			url := tiles[h][w].ToMapURL(datatype, "png")
			base, err := internal.FetchImage(url) // download base map
			if err != nil {
				return nil, err
			}
			imgs[h][w] = base
		}
	}
	img := ConcatImages(imgs)
	return CropImage(img, tc, rect), nil
}

func FetchBorderImage(gc GeoCoordinate, zoom uint, rect Rect) (image.Image, error) {
	tc := gc.GetTileCoordinate(zoom)
	tiles := makeTiles(tc, rect)
	imgs := make([][]image.Image, len(tiles))
	for h := range imgs {
		imgs[h] = make([]image.Image, len(tiles[0]))
	}
	for h := range imgs {
		for w := range imgs[h] {
			url := tiles[h][w].ToBorderMapURL("png")
			border, err := internal.FetchImage(url) // download prefectural border map
			if err != nil {
				return nil, err
			}
			imgs[h][w] = border
		}
	}
	img := ConcatImages(imgs)
	return CropImage(img, tc, rect), nil
}

func FetchJmaImage(gc GeoCoordinate, zoom uint, rect Rect, now time.Time, duration time.Duration) (image.Image, error) {
	tc := gc.GetTileCoordinate(zoom)
	tiles := makeTiles(tc, rect)
	imgs := make([][]image.Image, len(tiles))
	for h := range imgs {
		imgs[h] = make([]image.Image, len(tiles[0]))
	}
	for h := range imgs {
		for w := range imgs[h] {
			url, err := tiles[h][w].ToJmaURL(now, duration, "png")
			if err != nil {
				return nil, err
			}
			weather, err := internal.FetchImage(url) // downlaod weather map
			if err != nil {
				return nil, err
			}
			imgs[h][w] = weather
		}
	}
	img := ConcatImages(imgs)
	return CropImage(img, tc, rect), nil
}

func makeTiles(tc TileCoordinate, rect Rect) [][]Tile {
	m, M := CalcCorner(tc, float64(rect.H)/256, float64(rect.H)/256)
	min := m.GetTile()
	max := M.GetTile()
	nh, nw := calcCanvasSize(tc, min, max)
	tiles := make([][]Tile, nh)
	for h := range tiles {
		tiles[h] = make([]Tile, nw)
		for w := range tiles[h] {
			tiles[h][w] = Tile{tc.Zoom, min.X + uint(w), min.Y + uint(h)}
		}
	}
	return tiles
}

func calcCanvasSize(tc TileCoordinate, min Tile, max Tile) (int, int) {
	var nh, nw int
	if isApproxEq(math.Mod(tc.Y, 1), 0.5) {
		nh = int(max.Y - min.Y)
	} else {
		nh = int(max.Y - min.Y + 1)
	}
	if isApproxEq(math.Mod(tc.X, 1), 0.5) {
		nw = int(max.X - min.X)
	} else {
		nw = int(max.X - min.X + 1)
	}
	return nh, nw
}

func isApproxEq(x float64, y float64) bool {
	eps := 1e-5
	return math.Abs(x-y) < eps
}
