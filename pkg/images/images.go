package images

import (
	"image"
	"io"

	"golang.org/x/image/draw"
)

func ResizeImage(ext string, input io.Reader, height int, width int) ([]byte, error) {
	if ext == "png" {
		return ResizePNG(input, height, width)
	} else if ext == "jpg" || ext == "jpeg" {
		return ResizeJPEG(input, height, width)
	}

	return nil, nil
}

func resizeImageIfNecessary(src image.Image, height int, width int) *image.RGBA {
	if src.Bounds().Max.X > height || src.Bounds().Max.Y > width {
		// Set the expected size we want
		dst := image.NewRGBA(image.Rect(0, 0, height, width))

		// Resize image
		draw.ApproxBiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

		return dst
	}

	return nil
}
