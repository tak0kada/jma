package jma

import (
	"errors"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
)

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

func CropImage(img image.Image, center TileCoordinate, rect Rect) image.Image {
	return imaging.CropCenter(img, int(rect.W), int(rect.H))
}

func Decolor(img image.Image) image.Image {
	return imaging.Grayscale(img)
}