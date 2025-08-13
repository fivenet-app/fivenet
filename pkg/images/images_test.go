package images

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"testing"

	webp "github.com/HugoSmits86/nativewebp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			ext:         PNGExt,
			inputImage:  createTestImage(100, 100),
			height:      50,
			width:       50,
			expectError: false,
		},
		{
			name:        "Valid JPEG Resize",
			ext:         JPEGExt,
			inputImage:  createTestImage(100, 100),
			height:      50,
			width:       50,
			expectError: false,
		},
		{
			name:        "Unsupported Extension",
			ext:         GIFExt,
			inputImage:  createTestImage(100, 100),
			height:      50,
			width:       50,
			expectError: true,
			err:         ErrUnsupportedImageType,
		},
		{
			name:        "Zero Dimensions",
			ext:         PNGExt,
			inputImage:  createTestImage(100, 100),
			height:      0,
			width:       0,
			expectError: true,
			err:         ErrZeroResize,
		},
		{
			name:        "Invalid Image Decode",
			ext:         PNGExt,
			inputImage:  nil,
			height:      50,
			width:       50,
			expectError: true,
			err:         io.ErrUnexpectedEOF,
		},
		{
			name:        "Invalid Image Decode but JPEG",
			ext:         JPEGExt,
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
			switch {
			case tt.inputImage == nil:

			case tt.ext == PNGExt:
				err := png.Encode(&buf, tt.inputImage)
				require.NoError(t, err)

			case tt.ext == JPEGExt:
				err := jpeg.Encode(&buf, tt.inputImage, nil)
				require.NoError(t, err)
			}

			// Resize image
			output, err := ResizeImage(tt.ext, &buf, tt.height, tt.width)
			if tt.expectError {
				require.Error(t, err)
				if tt.err != nil {
					assert.Equal(t, tt.err, err)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)

				// Make sure the image was resized accordingly
				switch {
				case tt.inputImage == nil:

				case tt.ext == PNGExt:
					img, err := png.Decode(bytes.NewBuffer(output))
					require.NoError(t, err)
					assert.Equal(t, int(tt.height), img.Bounds().Dy())
					assert.Equal(t, int(tt.width), img.Bounds().Dx())

				case tt.ext == "jpeg":
					img, err := jpeg.Decode(bytes.NewBuffer(output))
					require.NoError(t, err)
					assert.Equal(t, int(tt.height), img.Bounds().Dy())
					assert.Equal(t, int(tt.width), img.Bounds().Dx())

				case tt.ext == "webp":
					img, err := webp.Decode(bytes.NewBuffer(output))
					require.NoError(t, err)
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
