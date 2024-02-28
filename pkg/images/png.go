package images

import (
	"bytes"
	"image/png"
	"io"
)

func ResizePNG(input io.Reader, height int, width int) (io.Reader, error) {
	// Decode the image (from PNG to image.Image)
	src, err := png.Decode(input)
	if err != nil {
		return nil, err
	}

	resizeImageIfNecessary(src, height, width)

	// Encode to output
	output := bytes.NewBuffer([]byte{})
	if err := png.Encode(output, src); err != nil {
		return nil, err
	}

	return output, nil
}
