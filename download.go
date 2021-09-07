package jma

import (
	"github.com/disintegration/imaging"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadImage(tile Tile, rect Rect, now time.Time, duration time.Duration, filepath string) error {
	img, err := FetchImage(tile, rect, now, duration)
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

func downloadImage(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fp, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fp, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
