package filestore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectContentType(t *testing.T) {
	t.Parallel()

	t.Run("Detects PNG", func(t *testing.T) {
		t.Parallel()

		pngHeader := []byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
			0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		}
		assert.Equal(t, "image/png", detectContentType(pngHeader))
	})

	t.Run("Plain text is detected", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "text/plain", detectContentType([]byte("not a recognized file format")))
	})

	t.Run("Opaque binary returns empty string", func(t *testing.T) {
		t.Parallel()

		data := []byte{0x00, 0xff, 0x81, 0x82, 0x00, 0x01, 0x02, 0x03}
		assert.Empty(t, detectContentType(data))
	})
}
