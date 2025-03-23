package images

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizeImage(t *testing.T) {
	tests := []struct {
		name        string
		ext         string
		inputImage  image.Image
		height      uint
		width       uint
		expectError bool
		err         error
	}{
		{
			name:        "Valid PNG Resize",
			ext:         "png",
			inputImage:  createTestImage(100, 100),
			height:      50,
			width:       50,
			expectError: false,
		},
		{
			name:        "Valid JPEG Resize",
			ext:         "jpeg",
			inputImage:  createTestImage(100, 100),
			height:      50,
			width:       50,
			expectError: false,
		},
		{
			name:        "Unsupported Extension",
			ext:         "gif",
			inputImage:  createTestImage(100, 100),
			height:      50,
			width:       50,
			expectError: true,
			err:         ErrUnsupportedImageType,
		},
		{
			name:        "Zero Dimensions",
			ext:         "png",
			inputImage:  createTestImage(100, 100),
			height:      0,
			width:       0,
			expectError: true,
			err:         ErrZeroResize,
		},
		{
			name:        "Invalid Image Decode",
			ext:         "png",
			inputImage:  nil,
			height:      50,
			width:       50,
			expectError: true,
			err:         io.ErrUnexpectedEOF,
		},
		{
			name:        "Invalid Image Decode but JPEG",
			ext:         "jpeg",
			inputImage:  nil,
			height:      50,
			width:       50,
			expectError: true,
			err:         io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode the test image to a buffer
			var buf bytes.Buffer
			if tt.inputImage == nil {
			} else if tt.ext == "png" {
				err := png.Encode(&buf, tt.inputImage)
				assert.NoError(t, err)
			} else if tt.ext == "jpeg" {
				err := jpeg.Encode(&buf, tt.inputImage, nil)
				assert.NoError(t, err)
			}

			// Resize image
			output, err := ResizeImage(tt.ext, &buf, tt.height, tt.width)
			if tt.expectError {
				assert.Error(t, err)
				if tt.err != nil {
					assert.Equal(t, tt.err, err)
				}
				assert.Nil(t, output)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)

				// Make sure the image was resized accordingly
				if tt.inputImage == nil {
				} else if tt.ext == "png" {
					img, err := png.Decode(bytes.NewBuffer(output))
					assert.NoError(t, err)
					assert.Equal(t, int(tt.height), img.Bounds().Dy())
					assert.Equal(t, int(tt.width), img.Bounds().Dx())
				} else if tt.ext == "jpeg" {
					img, err := jpeg.Decode(bytes.NewBuffer(output))
					assert.NoError(t, err)
					assert.Equal(t, int(tt.height), img.Bounds().Dy())
					assert.Equal(t, int(tt.width), img.Bounds().Dx())
				}
			}
		})
	}
}

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := range width {
		for y := range height {
			img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	return img
}
