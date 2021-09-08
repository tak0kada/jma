package jma

import (
	"errors"
	"github.com/disintegration/imaging"
	"github.com/tak0kada/jma/internal"
	"image"
	"image/color"
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
	img, _ := Overlay(decolor(base), weather, border)
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
	tiles := initTiles(gc.ToTile(zoom), rect)
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
	tc := gc.ToTileCoordinate(zoom)
	return clipImage(img, tc, rect), nil
}

func FetchBorderImage(gc GeoCoordinate, zoom uint, rect Rect) (image.Image, error) {
	tiles := initTiles(gc.ToTile(zoom), rect)
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
	tc := gc.ToTileCoordinate(zoom)
	return clipImage(img, tc, rect), nil
}

func FetchJmaImage(gc GeoCoordinate, zoom uint, rect Rect, now time.Time, duration time.Duration) (image.Image, error) {
	tiles := initTiles(gc.ToTile(zoom), rect)
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
	tc := gc.ToTileCoordinate(zoom)
	return clipImage(img, tc, rect), nil
}

func ConcatImages(imgs [][]image.Image) image.Image {
	nh := len(imgs)
	nw := len(imgs[0])
	dst := imaging.New(256*nh, 256*nw, color.RGBA{0, 0, 0, 0})
	for h := range imgs {
		for w := range imgs[h] {
			dst = imaging.Paste(dst, imgs[h][w], image.Pt(256*h, 256*w))
		}
	}
	return dst
}

func Overlay(bottom image.Image, middle image.Image, top image.Image) (image.Image, error) {
	eqsize := func(left image.Image, right image.Image) bool {
		return left.Bounds().Dx() == right.Bounds().Dx() && left.Bounds().Dy() == right.Bounds().Dy()
	}
	if !(eqsize(bottom, middle) && eqsize(middle, top)) {
		return nil, errors.New("error: size of input images are not consistent")
	}
	opacity := 1.0
	dst := imaging.New(top.Bounds().Dx(), top.Bounds().Dy(), color.RGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, bottom, image.Pt(0, 0))
	dst = imaging.OverlayCenter(dst, middle, opacity)
	dst = imaging.OverlayCenter(dst, top, opacity)
	return dst, nil
}

func clipImage(img image.Image, tc TileCoordinate, rect Rect) image.Image {
	return imaging.CropCenter(img, int(rect.W), int(rect.H))
}

func initTiles(tile Tile, rect Rect) [][]Tile {
	nw, nh := calcCanvasSize(rect)
	tiles := make([][]Tile, nh)
	for h := range tiles {
		tiles[h] = make([]Tile, nw)
		for w := range tiles[h] {
			tiles[h][w] = Tile{tile.Zoom, tile.X - (nh-1)/2 + uint(h), tile.Y - (nw-1)/2 + uint(w)}
		}
	}
	return tiles
}

func calcCanvasSize(rect Rect) (uint, uint) {
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

func decolor(img image.Image) image.Image {
	return imaging.Grayscale(img)
}
