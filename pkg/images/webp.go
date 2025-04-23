package images

import (
	"image"
	"io"

	webp "github.com/HugoSmits86/nativewebp"
)

type WebP struct {
	ImageType
}

func (g WebP) Decode(input io.Reader) (image.Image, error) {
	// Decode the image (from WebP to image.Image)
	src, err := webp.Decode(input)
	if err != nil {
		return nil, err
	}

	return src, err
}

func (g WebP) Encode(w io.Writer, m image.Image) error {
	// Encode the image (from image.Image to WebP)
	if err := webp.Encode(w, m, &webp.Options{
		UseExtendedFormat: false,
	}); err != nil {
		return err
	}

	return nil
}
