package checksum

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/hex"
)

func TestMD5ByPath(t *testing.T) {
	actual, err := MD5ByPath("../test_files/checksum/3b5d5c3712955042212316173ccf37be")
	if err != nil {
		t.Fatal(err)
	}
	expected,err := hex.DecodeString("3b5d5c3712955042212316173ccf37be")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)
}
