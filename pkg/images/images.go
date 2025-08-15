package images

import (
	"bytes"
	"errors"
	"image"
	"io"
	"math"
	"strings"

	"golang.org/x/image/draw"
)

const (
	PNGExt  = "png"
	GIFExt  = "gif"
	JPGExt  = "jpg"
	JPEGExt = "jpeg"
	WEBPExt = "webp"
)

var (
	ErrUnsupportedImageType  = errors.New("unsupported image type")
	ErrZeroDimensions        = errors.New("image has zero dimensions")
	ErrZeroResize            = errors.New("resize height and width are both zero")
	ErrDimensionsOutOfBounds = errors.New("image dimensions exceed uint size")
)

// ImageType defines an interface for image encoding and decoding.
type ImageType interface {
	Decode(input io.Reader) (image.Image, error)
	Encode(w io.Writer, m image.Image) error
}

// ResizeImage resizes an image to the specified height and width.
func ResizeImage(ext string, input io.Reader, height uint, width uint) ([]byte, error) {
	var imgType ImageType

	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	switch ext {
	case PNGExt:
		imgType = PNG{}
	case JPGExt, JPEGExt:
		imgType = JPEG{}
	case WEBPExt:
		imgType = WebP{}
	default:
		return nil, ErrUnsupportedImageType
	}

	// Decode image to image.Image
	img, err := imgType.Decode(input)
	if err != nil {
		return nil, err
	}

	dst, err := resizeImageIfNecessary(img, height, width)
	if err != nil {
		return nil, err
	}
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

// Resize condition is heavily inspired by <https://github.com/KononK/resize> code.
func resizeImageIfNecessary(src image.Image, height uint, width uint) (*image.RGBA, error) {
	// Source image has no or "negative" pixels
	if src.Bounds().Dx() <= 0 || src.Bounds().Dy() <= 0 {
		return nil, ErrZeroDimensions
	}

	// Make sure the dimensions are within the bounds of uint
	if src.Bounds().Dx() > math.MaxInt || src.Bounds().Dy() > math.MaxInt {
		return nil, ErrDimensionsOutOfBounds
	}

	if width == 0 && height == 0 {
		return nil, ErrZeroResize
	}

	scaleX, scaleY := calcScaleFactors(
		width,
		height,
		float64(src.Bounds().Dx()),
		float64(src.Bounds().Dy()),
	)
	if width == 0 {
		width = uint(0.7 + float64(src.Bounds().Dx())/scaleX)
	}
	if height == 0 {
		height = uint(0.7 + float64(src.Bounds().Dy())/scaleY)
	}

	// Ensure the width and height are within bounds of int for the RGBA image creation.
	if width > math.MaxInt || height > math.MaxInt {
		return nil, ErrDimensionsOutOfBounds
	}

	// Nothing to do, return src image
	if width == uint(src.Bounds().Dx()) && height == uint(src.Bounds().Dy()) {
		return nil, nil
	}

	// Create the destination image with the expected size we want
	//nolint:gosec // Above width and height are checked to be within bounds of int.
	dst := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// Resize image go's draw builtin bilinear interpolator
	draw.ApproxBiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return dst, nil
}

// Calculates scaling factors using old and new image dimensions.
func calcScaleFactors(
	width uint,
	height uint,
	oldWidth float64,
	oldHeight float64,
) (float64, float64) {
	var scaleX float64
	var scaleY float64

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

	return scaleX, scaleY
}
