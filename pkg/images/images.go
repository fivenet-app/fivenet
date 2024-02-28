package images

import (
	"image"
	"io"

	"github.com/h2non/filetype/types"
	"golang.org/x/image/draw"
)

func ResizeImage(iType types.Type, input io.Reader, height int, width int) (io.Reader, error) {
	if iType.Extension == "png" {
		return ResizePNG(input, height, width)
	} else if iType.Extension == "jpg" || iType.Extension == "jpeg" {
		return ResizeJPEG(input, height, width)
	}

	return input, nil
}

func resizeImageIfNecessary(src image.Image, height int, width int) {
	if src.Bounds().Max.X > height || src.Bounds().Max.Y > width {
		// Set the expected size we want
		dst := image.NewRGBA(image.Rect(0, 0, height, width))

		// Resize image
		draw.ApproxBiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	}
}
