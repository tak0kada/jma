package jma

import (
	"github.com/disintegration/imaging"
	"io"
	"net/http"
	"os"
	"time"
)

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
