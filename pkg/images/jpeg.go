package images

import (
	"bytes"
	"image/jpeg"
	"io"
)

func ResizeJPEG(input io.Reader, height int, width int) ([]byte, error) {
	// Decode the image (from JPG to image.Image):
	src, err := jpeg.Decode(input)
	if err != nil {
		return nil, err
	}

	dst := resizeImageIfNecessary(src, height, width)
	if dst == nil {
		return nil, nil
	}

	// Encode to output
	output := new(bytes.Buffer)
	if err := jpeg.Encode(output, dst, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
