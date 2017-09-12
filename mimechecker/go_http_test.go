package mimechecker

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHttp_ByPath(t *testing.T) {
	mimeMap := map[string]string{
		"audio.mp3": "audio/mpeg",
		"image-gif-fake-audio.mp3": "image/gif",
		"audio-no-extension": "audio/mpeg",
		"image.png": "image/png",
		"image.jpg": "image/jpeg",
		"pdf.pdf": "application/pdf",
		"text.txt": "",
		"text-with-bom.txt": "text/plain; charset=utf-8",
		"yaml.yml": "",
	}
	checker := NewGoHttp()
	for filename, expectedMime := range mimeMap {
		result,_ := checker.TypeByFile("../example/" + filename)
		assert.Equal(t, expectedMime, result, "file: %s", filename)
	}
}
