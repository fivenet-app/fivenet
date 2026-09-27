package filestore

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFileName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "my_file.txt", SanitizeFileName("../my file.txt"))
	assert.Equal(t, "report.tar.gz", SanitizeFileName("report.tar.gz"))
	assert.Empty(t, SanitizeFileName(""))

	name := strings.Repeat("a", 250) + ".pdf"
	got := SanitizeFileName(name)
	assert.Len(t, got, 244)
	assert.Equal(t, ".pdf", got[240:])
}

func TestSniff(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "IMAGE/PNG", sniff("IMAGE/PNG", "file.txt"))
	assert.Equal(t, "application/pdf", sniff("", "file.pdf"))
	assert.Equal(t, ApplicationOctetStreamMIMEType, sniff("", "file.unknown"))
}

func TestCountingReader(t *testing.T) {
	t.Parallel()

	reader := &countingReader{Reader: bytes.NewBufferString("hello")}
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), data)
	assert.EqualValues(t, 5, reader.n)
}
