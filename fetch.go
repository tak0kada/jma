package jma

import (
	"bytes"
	"image"
	"io"
	"net/http"
	"time"
)

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
