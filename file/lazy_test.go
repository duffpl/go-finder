package file

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewLazyByPath(t *testing.T) {
	t.Run("Callbacks", func(t *testing.T) {
		t.Run("Checksum", func(t *testing.T) {
			expected := []byte("1234")
			info,_ := NewLazyFileInfoExByPath("../test_files/checksum/3b5d5c3712955042212316173ccf37be", func(path string) ([]byte, error) {
				return expected, nil
			}, func(path string) (string, error) {
				return "", nil
			})
			actual, _ := info.Checksum()
			assert.Equal(t, expected, actual)
		})
		t.Run("Mime", func(t *testing.T) {
			expected := ("mime123")
			info,_ := NewLazyFileInfoExByPath("../test_files/checksum/3b5d5c3712955042212316173ccf37be", func(path string) ([]byte, error) {
				return nil,nil
			}, func(path string) (string, error) {
				return expected, nil
			})
			actual, _ := info.Mime()
			assert.Equal(t, expected, actual)
		})
	})
}