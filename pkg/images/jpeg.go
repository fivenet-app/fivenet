package images

import (
	"image"
	"image/jpeg"
	"io"
)

type JPEG struct {
	ImageType
}

func (g JPEG) Decode(input io.Reader) (image.Image, error) {
	// Decode the image (from JPEG to image.Image)
	src, err := jpeg.Decode(input)
	if err != nil {
		return nil, err
	}

	return src, err
}

func (g JPEG) Encode(w io.Writer, m image.Image) error {
	// Encode the image (from image.Image to JPEG)
	if err := jpeg.Encode(w, m, &jpeg.Options{Quality: 90}); err != nil {
		return err
	}

	return nil
}
