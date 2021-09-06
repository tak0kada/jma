package jma

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

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

func FetchImage(url string) (io.Reader, error) {
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
