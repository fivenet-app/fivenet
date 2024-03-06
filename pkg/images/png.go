package images

import (
	"image"
	"image/png"
	"io"
)

type PNG struct {
	ImageType
}

func (g PNG) Decode(input io.Reader) (image.Image, error) {
	// Decode the image (from PNG to image.Image)
	src, err := png.Decode(input)
	if err != nil {
		return nil, err
	}

	return src, err
}

func (g PNG) Encode(w io.Writer, m image.Image) error {
	if err := png.Encode(w, m); err != nil {
		return err
	}

	return nil
}
