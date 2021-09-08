package internal

import (
	"bytes"
	"image"
	"io"
	"net/http"
	"os"
)

func FetchImage(url string) (image.Image, error) {
	reader, err := FetchImageReader(url)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, err
}

func FetchImageReader(url string) (io.Reader, error) {
	data, err := FetchImageByte(url)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func FetchImageByte(url string) ([]byte, error) {
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

func DownloadImage(url string, filepath string) error {
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
