package jma

import (
	"github.com/disintegration/imaging"
	"time"
)

func DownloadImage(gc GeoCoordinate, zoom uint, rect Rect, now time.Time, duration time.Duration, filepath string) error {
	img, err := FetchImage(gc, zoom, rect, now, duration)
	if err != nil {
		return err
	}
	return imaging.Save(img, filepath)
}

func DownloadImageTile(tile Tile, now time.Time, duration time.Duration, filepath string) error {
	img, err := FetchImageTile(tile, now, duration)
	if err != nil {
		return err
	}
	return imaging.Save(img, filepath)
}
