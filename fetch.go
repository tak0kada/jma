package jma

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"io"
	"net/http"
	"time"
)

func FetchImage(tile Tile, rect Rect, now time.Time, duration time.Duration) (image.Image, error) {
	tiles := initTiles(tile, rect)
	imgs := make([][]image.Image, len(tiles))
	for h := range imgs {
		imgs[h] = make([]image.Image, len(tiles[0]))
	}

	var err error
	for h := range imgs {
		for w := range imgs[h] {
			imgs[h][w], err = FetchImageTile(tiles[h][w], now, duration)
			if err != nil {
				return nil, err
			}
		}
	}

	img := concatImages(imgs)
	return imaging.CropCenter(img, int(rect.W), int(rect.H)), nil
}

func concatImages(imgs [][]image.Image) image.Image {
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

func FetchImageTile(tile Tile, now time.Time, duration time.Duration) (image.Image, error) {
	url := tile.ToMapURL("pale", "png")
	base, err := fetchImage(url) // download base map
	if err != nil {
		return nil, err
	}

	url = tile.ToBorderMapURL("png")
	border, err := fetchImage(url) // download prefectural border map
	if err != nil {
		return nil, err
	}

	url, err = tile.ToJmaURL(now, duration, "png")
	weather, err := fetchImage(url) // downlaod weather map
	if err != nil {
		return nil, err
	}

	return overlay(decolor(base), border, weather), nil
}

func overlay(bottom image.Image, middle image.Image, top image.Image) image.Image {
	dst := imaging.New(256, 256, color.RGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, bottom, image.Pt(0, 0))
	dst = imaging.OverlayCenter(dst, middle, 1.0)
	dst = imaging.OverlayCenter(dst, top, 1.0)
	return dst
}

func decolor(img image.Image) image.Image {
	return imaging.Grayscale(img)
}

func fetchImage(url string) (image.Image, error) {
	reader, err := fetchImageReader(url)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, err
}

func fetchImageReader(url string) (io.Reader, error) {
	data, err := fetchImageByte(url)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func fetchImageByte(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
