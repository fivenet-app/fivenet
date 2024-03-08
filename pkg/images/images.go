package images

import (
	"bytes"
	"image"
	"io"

	"golang.org/x/image/draw"
)

type ImageType interface {
	Decode(input io.Reader) (image.Image, error)
	Encode(w io.Writer, m image.Image) error
}

type ResizeTypeFn = func(input io.Reader, height uint, width uint) ([]byte, error)

func ResizeImage(ext string, input io.Reader, height uint, width uint) ([]byte, error) {
	var imgType ImageType

	switch ext {
	case "png":
		imgType = PNG{}
	case "jpg":
		fallthrough
	case "jpeg":
		imgType = JPEG{}
	default:
		return nil, nil
	}

	// Decode image to image.Image
	img, err := imgType.Decode(input)
	if err != nil {
		return nil, err
	}

	dst := resizeImageIfNecessary(img, height, width)
	if dst == nil {
		return nil, nil
	}

	output := new(bytes.Buffer)
	// Encode image to byte slices
	if err := imgType.Encode(output, dst); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// Resize logic is heavily inspired by https://github.com/KononK/resize code
func resizeImageIfNecessary(src image.Image, height uint, width uint) *image.RGBA {
	// Source image has no or "negative" pixels
	if src.Bounds().Dx() <= 0 || src.Bounds().Dy() <= 0 {
		return nil
	}

	scaleX, scaleY := calcScaleFactors(width, height, float64(src.Bounds().Dx()), float64(src.Bounds().Dy()))
	if width == 0 {
		width = uint(0.7 + float64(src.Bounds().Dx())/scaleX)
	}
	if height == 0 {
		height = uint(0.7 + float64(src.Bounds().Dy())/scaleY)
	}

	// Nothing to do, return src image
	if int(width) == src.Bounds().Dx() && int(height) == src.Bounds().Dy() {
		return nil
	}

	// Create the destination image with the expected size we want
	dst := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// Resize image go's draw builtin bilinear interpolator
	draw.ApproxBiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return dst
}

// Calculates scaling factors using old and new image dimensions.
func calcScaleFactors(width uint, height uint, oldWidth float64, oldHeight float64) (scaleX float64, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}

	return
}
