package jma

import (
	"github.com/tak0kada/jma/internal"
	"image"
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
	img, _ := Overlay(Decolor(base), weather, border)
	return img, nil
}

func FetchImageTile(tile Tile, now time.Time, duration time.Duration) (image.Image, error) {
	tc := TileCoordinate{
		Zoom: tile.Zoom,
		X:    float64(tile.X),
		Y:    float64(tile.Y),
	}
	return FetchImage(tc.ToGeoCoordinate(), tc.Zoom, Rect{256, 256}, now, duration)
}

func FetchMapImage(gc GeoCoordinate, zoom uint, rect Rect, datatype string) (image.Image, error) {
	tc := gc.ToTileCoordinate(zoom)
	tiles := initTiles(tc, rect)
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
	tc := gc.ToTileCoordinate(zoom)
	tiles := initTiles(tc, rect)
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
	tc := gc.ToTileCoordinate(zoom)
	tiles := initTiles(tc, rect)
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

func initTiles(tc TileCoordinate, rect Rect) [][]Tile {
	tile := tc.ToTile()
	nw, nh := calcCanvasSize(tc, rect)
	tiles := make([][]Tile, nh)
	for h := range tiles {
		tiles[h] = make([]Tile, nw)
		for w := range tiles[h] {
			tiles[h][w] = Tile{tile.Zoom, tile.X - (nw-1)/2 + uint(w), tile.Y - (nh-1)/2 + uint(h)}
		}
	}
	return tiles
}

func calcCanvasSize(tc TileCoordinate, rect Rect) (uint, uint) {
	var w, h uint
	if rect.W/2-128 > 0 {
		w = 2*((rect.W/2-128)/256+1) + 1
	} else {
		w = 1
	}
	if rect.H/2-128 > 0 {
		h = 2*((rect.H/2-128)/256+1) + 1
	} else {
		h = 1
	}
	return w, h
}
