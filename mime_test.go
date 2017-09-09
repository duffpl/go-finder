package finder

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/duffpl/go-finder/mimechecker"
	"github.com/rakyll/magicmime"
	"log"
)

func TestStdMimeChecker(t *testing.T) {
	mimeMap := map[string]string{
		"audio.mp3": "audio/mpeg",
		"image-gif-fake-audio.mp3": "image/gif",
		"audio-no-extension": "audio/mpeg",
		"image.png": "image/png",
		"image.jpg": "image/jpeg",
		"pdf.pdf": "application/pdf",
		"text.txt": "text/plain",
		"text-with-bom.txt": "text/plain",
		"yaml.yml": "text/plain",
	}
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	defer magicmime.Close()
	stdMimeChecker = mimechecker.NewMagicMime()
	for filename, expectedMime := range mimeMap {
		result,_ := stdMimeChecker.ByPath("./example/" + filename)
		assert.Equal(t, expectedMime, result, "file: %s", filename)
	}
}
