package jma

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

func init() {
	validate.RegisterStructValidation(tileXYValidation, Tile{})
	validate.RegisterStructValidation(tileCoordinateXYValidation, TileCoordinate{})
}
